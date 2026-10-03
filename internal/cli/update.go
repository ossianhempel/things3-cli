package cli

import (
	"fmt"
	"strings"

	"github.com/ossianhempel/things3-cli/internal/db"
	"github.com/ossianhempel/things3-cli/internal/things"
	"github.com/spf13/cobra"
)

// NewUpdateCommand builds the update subcommand.
func NewUpdateCommand(app *App) *cobra.Command {
	opts := things.UpdateOptions{}
	repeatOpts := RepeatOptions{}
	var dbPath string
	var allowUnsafeTitle bool
	var noVerify bool
	var allowNonToday bool
	var yes bool
	var repeatJSON bool
	queryOpts := TaskQueryOptions{
		Status: "incomplete",
		Limit:  200,
	}

	cmd := &cobra.Command{
		Use:   "update [OPTIONS...] [--] [-|TITLE]",
		Short: "Update an existing todo",
		RunE: func(cmd *cobra.Command, args []string) error {
			rawInput, err := readInput(app.In, args)
			if err != nil {
				return err
			}

			repeatSpec, err := parseRepeatSpec(cmd, repeatOpts)
			if err != nil {
				return err
			}
			if repeatSpec.Enabled && strings.TrimSpace(opts.ID) == "" {
				return fmt.Errorf("Error: repeating updates require --id")
			}
			if repeatSpec.Enabled && (opts.Completed || opts.Canceled) {
				return fmt.Errorf("Error: --completed/--canceled cannot be combined with repeat changes because repeat templates must remain incomplete")
			}
			if repeatSpec.Enabled && opts.Duplicate {
				return fmt.Errorf("Error: --duplicate cannot be combined with repeat changes because the repeat rule would apply to the source template, not the duplicate")
			}
			if repeatJSON && !repeatSpec.Enabled {
				return fmt.Errorf("Error: --json is currently supported only for repeat updates")
			}
			if repeatSpec.Enabled && !repeatSpec.Clear && opts.Deadline != "" && repeatSpec.Spec.DeadlineOffset != nil {
				return fmt.Errorf("Error: --deadline cannot be combined with --repeat-deadline")
			}
			if repeatSpec.Enabled && !repeatSpec.Clear && (cmd.Flags().Changed("when") || opts.Later) {
				return fmt.Errorf("Error: --when/--later cannot be combined with --repeat because repeat activation controls scheduling; apply the repeat first, then update the template separately")
			}
			title := extractTitle(rawInput, "")
			if err := guardUnsafeTitle(title, allowUnsafeTitle); err != nil {
				return err
			}
			if err := validateWhenInput(opts.When); err != nil {
				return err
			}
			opts.NotesSet = cmd.Flags().Changed("notes")
			opts.WhenSet = cmd.Flags().Changed("when")
			opts.DeadlineSet = cmd.Flags().Changed("deadline")
			opts.TagsSet = cmd.Flags().Changed("tags")
			verifyWhen := resolveWhenValue(opts.When, opts.Later)
			verifyWhenEnabled := verifyWhen != "" && !noVerify && !app.DryRun
			guardEvening := strings.EqualFold(verifyWhen, "evening") && !allowNonToday
			ensureAuth := func() error {
				token, err := resolveAuthToken(app, opts.AuthToken)
				if err != nil {
					return err
				}
				opts.AuthToken = token
				return nil
			}

			queryOpts.HasURLSet = cmd.Flags().Changed("has-url")
			changedStatus := cmd.Flags().Changed("status")
			if strings.TrimSpace(opts.ID) != "" && hasExplicitSelector(map[string]bool{"status": changedStatus}, queryOpts) {
				return fmt.Errorf("Error: use either --id or query filters")
			}

			if strings.TrimSpace(opts.ID) == "" {
				if hasChecklistStatusChanges(opts) {
					return fmt.Errorf("Error: checklist item status updates require --id")
				}
				if !hasExplicitSelector(map[string]bool{"status": changedStatus}, queryOpts) {
					if err := ensureAuth(); err != nil {
						return err
					}
					url, err := things.BuildUpdateURL(opts, rawInput)
					if err != nil {
						return err
					}
					return openURL(app, url)
				}
				store, _, err := db.OpenDefault(dbPath)
				if err != nil {
					return formatDBError(err)
				}
				defer store.Close()

				tasks, err := fetchTasks(store, store.Tasks, queryOpts, false, []int{db.TaskTypeTodo})
				if err != nil {
					return formatDBError(err)
				}
				if len(tasks) == 0 {
					return fmt.Errorf("Error: no tasks matched")
				}
				if rawInput != "" && len(tasks) > 1 {
					return fmt.Errorf("Error: bulk update does not accept input (use --id or refine the query)")
				}
				if app.DryRun {
					return previewTasks(app.Out, tasks)
				}
				if guardEvening {
					for _, task := range tasks {
						if err := validateEveningTask(task, allowNonToday); err != nil {
							return err
						}
					}
				}
				if verifyWhen != "" {
					for _, task := range tasks {
						if task.Repeating {
							return fmt.Errorf("Error: cannot update when for repeating todos (id %s)", task.UUID)
						}
					}
				}
				if len(tasks) > 1 && !yes {
					return fmt.Errorf("Error: %d tasks matched (rerun with --yes to apply)", len(tasks))
				}
				if err := ensureAuth(); err != nil {
					return err
				}

				entry := ActionEntry{
					Type:  ActionUpdate,
					Items: make([]ActionItem, 0, len(tasks)),
				}
				for _, task := range tasks {
					entry.Items = append(entry.Items, taskToActionItem(task))
				}
				if err := appendAction(entry); err != nil {
					fmt.Fprintf(app.Err, "Warning: failed to write action log: %v\n", err)
				}

				for _, task := range tasks {
					opts.ID = task.UUID
					url, err := things.BuildUpdateURL(opts, rawInput)
					if err != nil {
						return err
					}
					if err := openURL(app, url); err != nil {
						return err
					}
					if verifyWhenEnabled {
						if err := verifyWhenApplied(store, task.UUID, verifyWhen); err != nil {
							return err
						}
					}
				}
				return nil
			}

			hasChanges := hasTodoUpdateChanges(opts, rawInput)
			if hasChecklistStatusChanges(opts) {
				if repeatSpec.Enabled {
					return fmt.Errorf("Error: checklist item status updates cannot be combined with --repeat")
				}
				if hasChangesWithoutChecklistStatus(opts, rawInput) {
					return fmt.Errorf("Error: checklist item status updates cannot be combined with other update fields")
				}
				if err := ensureAuth(); err != nil {
					return err
				}
				store, _, err := db.OpenDefault(dbPath)
				if err != nil {
					return formatDBError(err)
				}
				defer store.Close()
				task, err := store.TaskByID(opts.ID)
				if err != nil {
					return formatDBError(err)
				}
				checklist := make([]things.ChecklistItemState, 0, len(task.Checklist))
				for _, item := range task.Checklist {
					checklist = append(checklist, things.ChecklistItemState{
						Title:  item.Title,
						Status: item.Status,
					})
				}
				url, err := things.BuildChecklistStatusUpdateURL(opts, checklist)
				if err != nil {
					return err
				}
				if !app.DryRun {
					entry := ActionEntry{
						Type:  ActionUpdate,
						Items: []ActionItem{taskToActionItem(*task)},
					}
					if err := appendAction(entry); err != nil {
						fmt.Fprintf(app.Err, "Warning: failed to write action log: %v\n", err)
					}
				}
				return openURL(app, url)
			}
			if !repeatSpec.Enabled {
				if err := ensureAuth(); err != nil {
					return err
				}
				var verifyStore *db.Store
				if verifyWhenEnabled {
					verifyStore, err = openVerifyStore(app, dbPath)
					if err != nil {
						return err
					}
					if verifyStore != nil {
						defer verifyStore.Close()
						if task, err := verifyStore.TaskByID(opts.ID); err == nil && task.Repeating {
							return fmt.Errorf("Error: cannot update when for repeating todos (id %s)", opts.ID)
						}
					}
				}

				url, err := things.BuildUpdateURL(opts, rawInput)
				if err != nil {
					return err
				}
				if app.DryRun {
					return openURL(app, url)
				}
				var logStore *db.Store
				if guardEvening {
					store, err := openVerifyStore(app, dbPath)
					if err != nil {
						return err
					}
					if store != nil {
						logStore = store
						if task, err := store.TaskByID(opts.ID); err == nil {
							if err := validateEveningTask(*task, allowNonToday); err != nil {
								store.Close()
								return err
							}
						}
					}
				}

				if logStore == nil {
					logStore, _, err = db.OpenDefault(dbPath)
					if err != nil {
						logStore = nil
					}
				}
				if logStore != nil {
					if task, err := logStore.TaskByID(opts.ID); err == nil {
						entry := ActionEntry{
							Type:  ActionUpdate,
							Items: []ActionItem{taskToActionItem(*task)},
						}
						if err := appendAction(entry); err != nil {
							fmt.Fprintf(app.Err, "Warning: failed to write action log: %v\n", err)
						}
					}
					logStore.Close()
				}

				if err := openURL(app, url); err != nil {
					return err
				}
				if verifyWhenEnabled && verifyStore != nil {
					if err := verifyWhenApplied(verifyStore, opts.ID, verifyWhen); err != nil {
						return err
					}
				}
				return nil
			}

			plan, err := prepareRepeatUpdate(dbPath, opts.ID, repeatSpec, app.DryRun)
			if err != nil {
				return formatDBError(err)
			}
			defer plan.store.Close()
			if hasChanges && plan.usedTemplate {
				return fmt.Errorf("Error: combined ordinary and repeat updates cannot target a generated occurrence; rerun against template UUID %s", plan.targetID)
			}

			if hasChanges {
				if err := ensureAuth(); err != nil {
					return err
				}
				var verifyStore *db.Store
				if verifyWhenEnabled {
					verifyStore, err = openVerifyStore(app, dbPath)
					if err != nil {
						return err
					}
					if verifyStore != nil {
						defer verifyStore.Close()
						if task, err := verifyStore.TaskByID(opts.ID); err == nil && task.Repeating {
							return fmt.Errorf("Error: cannot update when for repeating todos (id %s)", opts.ID)
						}
					}
				}

				url, err := things.BuildUpdateURL(opts, rawInput)
				if err != nil {
					return err
				}
				if app.DryRun {
					if !repeatJSON {
						if err := openURL(app, url); err != nil {
							return err
						}
					}
					plan.result.Intent = &repeatIntent{URL: redactURLIntent(url)}
					plan.result.addStage(repeatStageURL, repeatStatusPlanned)
					plan.result.addStage(repeatStageDatabase, repeatStatusPlanned)
					plan.result.addStage(repeatStageVerification, repeatStatusPlanned)
					return renderRepeatResult(app.Out, plan.result, repeatJSON)
				}
				var logStore *db.Store
				if guardEvening {
					store, err := openVerifyStore(app, dbPath)
					if err != nil {
						return err
					}
					if store != nil {
						logStore = store
						if task, err := store.TaskByID(opts.ID); err == nil {
							if err := validateEveningTask(*task, allowNonToday); err != nil {
								store.Close()
								return err
							}
						}
					}
				}

				if logStore == nil {
					logStore, _, err = db.OpenDefault(dbPath)
					if err != nil {
						logStore = nil
					}
				}
				if logStore != nil {
					if task, err := logStore.TaskByID(opts.ID); err == nil {
						entry := ActionEntry{
							Type:  ActionUpdate,
							Items: []ActionItem{taskToActionItem(*task)},
						}
						if err := appendAction(entry); err != nil {
							fmt.Fprintf(app.Err, "Warning: failed to write action log: %v\n", err)
						}
					}
					logStore.Close()
				}

				if err := openURL(app, url); err != nil {
					return err
				}
				plan.result.Intent = &repeatIntent{URL: redactURLIntent(url)}
				plan.result.addStage(repeatStageURL, repeatStatusCompleted)
				if verifyWhenEnabled && verifyStore != nil {
					if err := verifyWhenApplied(verifyStore, opts.ID, verifyWhen); err != nil {
						return err
					}
				}
			} else if app.DryRun {
				plan.result.addStage(repeatStageDatabase, repeatStatusPlanned)
				plan.result.addStage(repeatStageVerification, repeatStatusPlanned)
				return renderRepeatResult(app.Out, plan.result, repeatJSON)
			}
			if plan.usedTemplate {
				fmt.Fprintf(app.Err, "Note: resolved repeating template %s for update\n", plan.targetID)
			}
			if err := plan.executeMutation(); err != nil {
				_ = renderRepeatResult(app.Out, plan.result, repeatJSON)
				return formatDBError(err)
			}
			return renderRepeatResult(app.Out, plan.result, repeatJSON)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&dbPath, "db", "d", "", "Path to Things database (overrides THINGSDB)")
	flags.StringVar(&dbPath, "database", "", "Alias for --db")
	flags.StringVar(&opts.AuthToken, "auth-token", "", "Things URL scheme authorization token")
	flags.StringVar(&opts.ID, "id", "", "ID of the todo to update")
	flags.StringVar(&opts.Notes, "notes", "", "Replace notes")
	flags.StringVar(&opts.PrependNotes, "prepend-notes", "", "Prepend to notes")
	flags.StringVar(&opts.AppendNotes, "append-notes", "", "Append to notes")
	flags.StringVar(&opts.When, "when", "", "When to schedule the todo (empty string clears the date)")
	flags.BoolVar(&opts.Later, "later", false, "Move the todo to Later")
	flags.StringVar(&opts.Deadline, "deadline", "", "Deadline for the todo (empty string clears it)")
	flags.StringVar(&opts.Tags, "tags", "", "Replace tags (empty string clears all tags)")
	flags.StringVar(&opts.AddTags, "add-tags", "", "Add tags")
	flags.BoolVar(&opts.Completed, "completed", false, "Mark the todo completed")
	flags.BoolVar(&opts.Canceled, "canceled", false, "Mark the todo canceled")
	flags.BoolVar(&opts.Canceled, "cancelled", false, "Mark the todo cancelled")
	flags.BoolVar(&opts.Reveal, "reveal", false, "Reveal the updated todo")
	flags.BoolVar(&opts.Duplicate, "duplicate", false, "Duplicate before updating")
	flags.StringVar(&opts.CompletionDate, "completion-date", "", "Completion date (ISO8601)")
	flags.StringVar(&opts.CreationDate, "creation-date", "", "Creation date (ISO8601)")
	flags.StringVar(&opts.Heading, "heading", "", "Heading within a project")
	flags.StringVar(&opts.List, "list", "", "Project or area to move to")
	flags.StringVar(&opts.ListID, "list-id", "", "Project or area ID to move to")
	flags.StringArrayVar(&opts.ChecklistItems, "checklist-item", nil, "Checklist item (repeatable)")
	flags.StringArrayVar(&opts.PrependChecklistItems, "prepend-checklist-item", nil, "Prepend checklist item (repeatable)")
	flags.StringArrayVar(&opts.AppendChecklistItems, "append-checklist-item", nil, "Append checklist item (repeatable)")
	flags.StringArrayVar(&opts.CompleteChecklistItems, "complete-checklist-item", nil, "Mark an existing checklist item complete by exact title (repeatable)")
	flags.StringArrayVar(&opts.IncompleteChecklistItems, "incomplete-checklist-item", nil, "Mark an existing checklist item incomplete by exact title (repeatable)")
	flags.BoolVar(&yes, "yes", false, "Confirm bulk update")
	flags.BoolVar(&allowUnsafeTitle, "allow-unsafe-title", false, "Allow titles that look like flag assignments")
	flags.BoolVar(&noVerify, "no-verify", false, "Skip verification of when updates against the Things database")
	flags.BoolVar(&allowNonToday, "allow-non-today", false, "Allow moving non-today tasks to This Evening")
	flags.BoolVar(&repeatJSON, "json", false, "Emit a structured result for repeat updates")
	addRepeatFlags(cmd, &repeatOpts, true)
	addTaskQueryFlags(cmd, &queryOpts, true, true)

	return cmd
}

func hasTodoUpdateChanges(opts things.UpdateOptions, rawInput string) bool {
	if strings.TrimSpace(rawInput) != "" {
		return true
	}
	if opts.Notes != "" || opts.NotesSet || opts.PrependNotes != "" || opts.AppendNotes != "" {
		return true
	}
	if opts.When != "" || opts.WhenSet || opts.Later {
		return true
	}
	if opts.Deadline != "" || opts.DeadlineSet {
		return true
	}
	if opts.Tags != "" || opts.TagsSet || opts.AddTags != "" {
		return true
	}
	if opts.Completed || opts.Canceled {
		return true
	}
	if opts.Reveal || opts.Duplicate {
		return true
	}
	if opts.CompletionDate != "" || opts.CreationDate != "" {
		return true
	}
	if opts.Heading != "" || opts.List != "" || opts.ListID != "" {
		return true
	}
	if len(opts.ChecklistItems) > 0 || len(opts.PrependChecklistItems) > 0 || len(opts.AppendChecklistItems) > 0 {
		return true
	}
	if hasChecklistStatusChanges(opts) {
		return true
	}
	return false
}

func hasChecklistStatusChanges(opts things.UpdateOptions) bool {
	return len(opts.CompleteChecklistItems) > 0 || len(opts.IncompleteChecklistItems) > 0
}

func hasChangesWithoutChecklistStatus(opts things.UpdateOptions, rawInput string) bool {
	withoutChecklistStatus := opts
	withoutChecklistStatus.CompleteChecklistItems = nil
	withoutChecklistStatus.IncompleteChecklistItems = nil
	return hasTodoUpdateChanges(withoutChecklistStatus, rawInput)
}
