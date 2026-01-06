package cli

import (
	"fmt"
	"time"

	"github.com/ossianhempel/things3-cli/internal/db"
	"github.com/ossianhempel/things3-cli/internal/repeat"
	"github.com/ossianhempel/things3-cli/internal/things"
	"github.com/spf13/cobra"
)

// NewAddCommand builds the add subcommand.
func NewAddCommand(app *App) *cobra.Command {
	opts := things.AddOptions{}
	repeatOpts := RepeatOptions{}
	var dbPath string

	cmd := &cobra.Command{
		Use:   "add [OPTIONS...] [--] [-|TITLE]",
		Short: "Add a new todo",
		RunE: func(cmd *cobra.Command, args []string) error {
			rawInput, err := readInput(app.In, args)
			if err != nil {
				return err
			}

			repeatSpec, err := parseRepeatSpec(cmd, repeatOpts)
			if err != nil {
				return err
			}
			title := extractTitle(rawInput, opts.TitlesRaw)

			if repeatSpec.Enabled {
				if repeatSpec.Clear {
					return fmt.Errorf("Error: --repeat-clear is only valid with update commands")
				}
				if opts.TitlesRaw != "" {
					return fmt.Errorf("Error: repeating add does not support --titles")
				}
				if opts.UseClipboard != "" {
					return fmt.Errorf("Error: repeating add does not support --use-clipboard")
				}
				if opts.ShowQuickEntry || title == "" {
					return fmt.Errorf("Error: repeating add requires an explicit title")
				}
			}

			url := things.BuildAddURL(opts, rawInput)
			if !repeatSpec.Enabled {
				return openURL(app, url)
			}
			if app.DryRun {
				if err := openURL(app, url); err != nil {
					return err
				}
				fmt.Fprintln(app.Err, "Note: --repeat is skipped in --dry-run mode.")
				return nil
			}

			ensureThingsLaunched(app)
			started := time.Now().Add(-2 * time.Second)
			if err := openURL(app, url); err != nil {
				return err
			}
			store, _, err := db.OpenDefaultWritable(dbPath)
			if err != nil {
				return formatDBError(err)
			}
			defer store.Close()

			taskID, err := waitForCreatedItem(store, title, db.TaskTypeTodo, started)
			if err != nil {
				return formatDBError(err)
			}
			update, err := repeat.BuildUpdate(repeatSpec.Spec)
			if err != nil {
				return err
			}
			if err := store.ApplyRepeatRule(taskID, update); err != nil {
				return formatDBError(err)
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&dbPath, "db", "d", "", "Path to Things database (overrides THINGSDB)")
	flags.StringVar(&dbPath, "database", "", "Alias for --db")
	flags.StringVar(&opts.When, "when", "", "When to schedule the todo")
	flags.StringVar(&opts.Deadline, "deadline", "", "Deadline for the todo")
	flags.BoolVar(&opts.Completed, "completed", false, "Mark the todo completed")
	flags.BoolVar(&opts.Canceled, "canceled", false, "Mark the todo canceled")
	flags.BoolVar(&opts.Canceled, "cancelled", false, "Mark the todo cancelled")
	flags.StringArrayVar(&opts.ChecklistItems, "checklist-item", nil, "Checklist item (repeatable)")
	flags.StringVar(&opts.CreationDate, "creation-date", "", "Creation date (ISO8601)")
	flags.StringVar(&opts.CompletionDate, "completion-date", "", "Completion date (ISO8601)")
	flags.StringVar(&opts.List, "list", "", "Project or area to add to")
	flags.StringVar(&opts.ListID, "list-id", "", "Project or area ID to add to")
	flags.StringVar(&opts.Heading, "heading", "", "Heading within a project")
	flags.BoolVar(&opts.Reveal, "reveal", false, "Reveal the newly created todo")
	flags.BoolVar(&opts.ShowQuickEntry, "show-quick-entry", false, "Show the quick entry dialog")
	flags.StringVar(&opts.Notes, "notes", "", "Notes for the todo")
	flags.StringVar(&opts.Tags, "tags", "", "Comma-separated tags")
	flags.StringVar(&opts.TitlesRaw, "titles", "", "Comma-separated titles for multiple todos")
	flags.StringVar(&opts.UseClipboard, "use-clipboard", "", "Use clipboard content")
	addRepeatFlags(cmd, &repeatOpts, false)

	return cmd
}
