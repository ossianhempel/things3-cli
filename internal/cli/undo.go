package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/ossianhempel/things3-cli/internal/db"
	"github.com/ossianhempel/things3-cli/internal/things"
	"github.com/spf13/cobra"
)

// NewUndoCommand builds the undo subcommand.
func NewUndoCommand(app *App) *cobra.Command {
	var authToken string
	var yes bool

	cmd := &cobra.Command{
		Use:   "undo",
		Short: "Undo the last bulk action",
		RunE: func(cmd *cobra.Command, args []string) error {
			entry, err := readLastAction()
			if err != nil {
				return fmt.Errorf("Error: %s", err)
			}
			if len(entry.Items) == 0 {
				return fmt.Errorf("Error: no actions logged")
			}

			if app.DryRun {
				fmt.Fprintf(app.Out, "Would undo %s for %d tasks\n", entry.Type, len(entry.Items))
				return previewActionItems(app.Out, entry.Items)
			}

			if len(entry.Items) > 1 && !yes {
				return fmt.Errorf("Error: %d tasks matched (rerun with --yes to apply)", len(entry.Items))
			}

			switch entry.Type {
			case ActionUpdate:
				token := strings.TrimSpace(authToken)
				if token == "" {
					token = authTokenFromEnv()
				}
				if token == "" {
					_, err := things.BuildUpdateURL(things.UpdateOptions{ID: "id"}, "")
					if err != nil {
						return err
					}
				}
				warnIncomplete := 0
				for _, item := range entry.Items {
					opts := things.UpdateOptions{
						AuthToken: token,
						ID:        item.UUID,
						Notes:     item.Notes,
						Tags:      strings.Join(item.Tags, ","),
						Deadline:  item.Deadline,
						Heading:   item.HeadingTitle,
					}
					when := whenFromActionItem(item)
					if when != "" {
						opts.When = when
					}
					if item.ProjectID != "" {
						opts.ListID = item.ProjectID
					} else if item.AreaID != "" {
						opts.ListID = item.AreaID
					}
					switch item.Status {
					case db.StatusCompleted:
						opts.Completed = true
					case db.StatusCanceled:
						opts.Canceled = true
					case db.StatusIncomplete:
						warnIncomplete++
					}
					url, err := things.BuildUpdateURL(opts, item.Title)
					if err != nil {
						return err
					}
					if err := openURL(app, url); err != nil {
						return err
					}
				}
				if warnIncomplete > 0 {
					fmt.Fprintln(app.Err, "Warning: Things URL scheme cannot un-complete tasks; some items may remain completed.")
				}
			case ActionTrash:
				for _, item := range entry.Items {
					opts := things.AddOptions{
						Notes:    item.Notes,
						Tags:     strings.Join(item.Tags, ","),
						Deadline: item.Deadline,
					}
					when := whenFromActionItem(item)
					if when != "" {
						opts.When = when
					}
					if item.ProjectID != "" {
						opts.ListID = item.ProjectID
					} else if item.AreaID != "" {
						opts.ListID = item.AreaID
					}
					if item.HeadingTitle != "" {
						opts.Heading = item.HeadingTitle
					}
					url := things.BuildAddURL(opts, item.Title)
					if err := openURL(app, url); err != nil {
						return err
					}
				}
				fmt.Fprintln(app.Err, "Warning: restored tasks are new items; trashed originals remain in Trash.")
			default:
				return fmt.Errorf("Error: unsupported action type %q", entry.Type)
			}

			if err := removeLastAction(); err != nil {
				fmt.Fprintf(app.Err, "Warning: failed to update action log: %v\n", err)
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&authToken, "auth-token", "", "Things URL scheme authorization token")
	flags.BoolVar(&yes, "yes", false, "Confirm undo for multiple tasks")

	return cmd
}

func whenFromActionItem(item ActionItem) string {
	if item.StartDate != "" {
		return item.StartDate
	}
	switch strings.ToLower(item.Start) {
	case "inbox":
		return "inbox"
	case "anytime":
		return "anytime"
	case "someday":
		return "someday"
	default:
		return ""
	}
}

func previewActionItems(out io.Writer, items []ActionItem) error {
	tasks := make([]db.Task, 0, len(items))
	for _, item := range items {
		tasks = append(tasks, db.Task{
			UUID:   item.UUID,
			Title:  item.Title,
			Status: item.Status,
		})
	}
	return previewTasks(out, tasks)
}
