package cli

import (
	"fmt"
	"strings"

	"github.com/ossianhempel/things3-cli/internal/db"
	"github.com/ossianhempel/things3-cli/internal/things"
	"github.com/spf13/cobra"
)

// NewDeleteCommand builds the delete subcommand.
func NewDeleteCommand(app *App) *cobra.Command {
	var dbPath string
	var id string
	var confirm string
	var yes bool
	opts := TaskQueryOptions{
		Status: "incomplete",
		Limit:  200,
	}

	cmd := &cobra.Command{
		Use:   "delete [OPTIONS...] [--] [-|TITLE]",
		Short: "Delete an existing todo",
		RunE: func(cmd *cobra.Command, args []string) error {
			rawInput, err := readInput(app.In, args)
			if err != nil {
				return err
			}

			changedStatus := cmd.Flags().Changed("status")
			opts.HasURLSet = cmd.Flags().Changed("has-url")
			hasTarget := strings.TrimSpace(id) != "" || strings.TrimSpace(rawInput) != ""
			if hasTarget {
				if hasExplicitSelector(map[string]bool{"status": changedStatus}, opts) {
					return fmt.Errorf("Error: use either --id or query filters")
				}

				target := deleteConfirmTarget(id, rawInput)
				if err := confirmDelete(app, "todo", target, confirm); err != nil {
					return err
				}

				script, err := things.BuildDeleteTodoScript(things.DeleteTodoOptions{ID: id}, rawInput)
				if err != nil {
					return err
				}
				return runScript(app, script)
			}

			if !hasExplicitSelector(map[string]bool{"status": changedStatus}, opts) {
				return fmt.Errorf("Error: refuse to delete without a selector (use --query/--search/--tag/etc)")
			}

			store, _, err := db.OpenDefault(dbPath)
			if err != nil {
				return formatDBError(err)
			}
			defer store.Close()

			tasks, err := fetchTasks(store, store.Tasks, opts, false, []int{db.TaskTypeTodo})
			if err != nil {
				return formatDBError(err)
			}
			if len(tasks) == 0 {
				return fmt.Errorf("Error: no tasks matched")
			}
			if app.DryRun {
				return previewTasks(app.Out, tasks)
			}
			if len(tasks) > 1 && !yes {
				return fmt.Errorf("Error: %d tasks matched (rerun with --yes to apply)", len(tasks))
			}

			entry := ActionEntry{
				Type:  ActionTrash,
				Items: make([]ActionItem, 0, len(tasks)),
			}
			for _, task := range tasks {
				entry.Items = append(entry.Items, taskToActionItem(task))
			}
			if err := appendAction(entry); err != nil {
				fmt.Fprintf(app.Err, "Warning: failed to write action log: %v\n", err)
			}

			ids := make([]string, 0, len(tasks))
			for _, task := range tasks {
				ids = append(ids, task.UUID)
			}
			script, err := things.BuildTrashScript(ids)
			if err != nil {
				return err
			}
			return runScript(app, script)
		},
	}

	cmd.Flags().StringVarP(&dbPath, "db", "d", "", "Path to Things database (overrides THINGSDB)")
	cmd.Flags().StringVar(&dbPath, "database", "", "Alias for --db")
	cmd.Flags().StringVar(&id, "id", "", "ID of the todo to delete")
	cmd.Flags().StringVar(&confirm, "confirm", "", "Confirm deletion by typing the todo ID or title")
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm bulk delete")
	addTaskQueryFlags(cmd, &opts, true, true)

	return cmd
}
