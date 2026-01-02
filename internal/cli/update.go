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
	var dbPath string
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

			if opts.AuthToken == "" {
				opts.AuthToken = authTokenFromEnv()
			}

			opts.AuthToken = strings.TrimSpace(opts.AuthToken)
			queryOpts.HasURLSet = cmd.Flags().Changed("has-url")
			changedStatus := cmd.Flags().Changed("status")
			if strings.TrimSpace(opts.ID) != "" && hasExplicitSelector(map[string]bool{"status": changedStatus}, queryOpts) {
				return fmt.Errorf("Error: use either --id or query filters")
			}

			if strings.TrimSpace(opts.ID) == "" {
				if !hasExplicitSelector(map[string]bool{"status": changedStatus}, queryOpts) {
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
				if len(tasks) > 1 && !yes {
					return fmt.Errorf("Error: %d tasks matched (rerun with --yes to apply)", len(tasks))
				}
				if opts.AuthToken == "" {
					_, err := things.BuildUpdateURL(things.UpdateOptions{ID: "id"}, "")
					if err != nil {
						return err
					}
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
				}
				return nil
			}

			url, err := things.BuildUpdateURL(opts, rawInput)
			if err != nil {
				return err
			}
			if app.DryRun {
				return openURL(app, url)
			}

			if opts.AuthToken == "" {
				_, err := things.BuildUpdateURL(things.UpdateOptions{ID: opts.ID}, "")
				if err != nil {
					return err
				}
			}

			store, _, err := db.OpenDefault(dbPath)
			if err == nil {
				if task, err := store.TaskByID(opts.ID); err == nil {
					entry := ActionEntry{
						Type:  ActionUpdate,
						Items: []ActionItem{taskToActionItem(*task)},
					}
					if err := appendAction(entry); err != nil {
						fmt.Fprintf(app.Err, "Warning: failed to write action log: %v\n", err)
					}
				}
				store.Close()
			}

			return openURL(app, url)
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
	addTaskQueryFlags(cmd, &queryOpts, true, true)

	return cmd
}
