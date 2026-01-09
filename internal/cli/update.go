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
			title := extractTitle(rawInput, "")
			if err := guardUnsafeTitle(title, allowUnsafeTitle); err != nil {
				return err
			}
			if err := validateWhenInput(opts.When); err != nil {
				return err
			}
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
					if err := openURL(app, url); err != nil {
						return err
					}
					if repeatSpec.Enabled {
						fmt.Fprintln(app.Err, "Note: --repeat is skipped in --dry-run mode.")
					}
					return nil
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
			} else if app.DryRun {
				fmt.Fprintf(app.Out, "Would update repeating rule for %s\n", opts.ID)
				return nil
			}
			if app.DryRun {
				fmt.Fprintln(app.Err, "Note: --repeat is skipped in --dry-run mode.")
				return nil
			}

			store, _, err := db.OpenDefaultWritable(dbPath)
			if err != nil {
				return formatDBError(err)
			}
			defer store.Close()

			targetID, usedTemplate, err := resolveRepeatTarget(store, opts.ID, db.TaskTypeTodo)
			if err != nil {
				return formatDBError(err)
			}
			if usedTemplate {
				fmt.Fprintf(app.Err, "Note: resolved repeating template %s for update\n", targetID)
			}
			if err := applyRepeatSpec(store, targetID, repeatSpec); err != nil {
				return formatDBError(err)
			}
			return nil
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
	flags.StringVar(&opts.When, "when", "", "When to schedule the todo")
	flags.BoolVar(&opts.Later, "later", false, "Move the todo to Later")
	flags.StringVar(&opts.Deadline, "deadline", "", "Deadline for the todo")
	flags.StringVar(&opts.Tags, "tags", "", "Replace tags")
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
	flags.BoolVar(&yes, "yes", false, "Confirm bulk update")
	flags.BoolVar(&allowUnsafeTitle, "allow-unsafe-title", false, "Allow titles that look like flag assignments")
	flags.BoolVar(&noVerify, "no-verify", false, "Skip verification of when updates against the Things database")
	flags.BoolVar(&allowNonToday, "allow-non-today", false, "Allow moving non-today tasks to This Evening")
	addRepeatFlags(cmd, &repeatOpts, true)
	addTaskQueryFlags(cmd, &queryOpts, true, true)

	return cmd
}

func hasTodoUpdateChanges(opts things.UpdateOptions, rawInput string) bool {
	if strings.TrimSpace(rawInput) != "" {
		return true
	}
	if opts.Notes != "" || opts.PrependNotes != "" || opts.AppendNotes != "" {
		return true
	}
	if opts.When != "" || opts.Later {
		return true
	}
	if opts.Deadline != "" {
		return true
	}
	if opts.Tags != "" || opts.AddTags != "" {
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
	return false
}
