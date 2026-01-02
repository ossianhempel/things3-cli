package cli

import (
	"github.com/ossianhempel/things3-cli/internal/db"
	"github.com/spf13/cobra"
)

// NewTasksCommand builds the tasks command.
func NewTasksCommand(app *App) *cobra.Command {
	var dbPath string
	opts := TaskQueryOptions{
		Status: "incomplete",
		Limit:  200,
	}
	var format string
	var selectRaw string
	var asJSON bool
	var noHeader bool

	cmd := &cobra.Command{
		Use:     "tasks",
		Short:   "List todos from the Things database",
		Aliases: []string{"todos"},
		RunE: func(cmd *cobra.Command, args []string) error {
			store, _, err := db.OpenDefault(dbPath)
			if err != nil {
				return formatDBError(err)
			}
			defer store.Close()

			opts.HasURLSet = cmd.Flags().Changed("has-url")
			outputOpts, err := resolveTaskOutputOptions(format, asJSON, selectRaw, noHeader)
			if err != nil {
				return err
			}
			tasks, err := fetchTasks(store, store.Tasks, opts, false, []int{db.TaskTypeTodo})
			if err != nil {
				return formatDBError(err)
			}
			return printTasks(app.Out, tasks, outputOpts)
		},
	}

	cmd.Flags().StringVarP(&dbPath, "db", "d", "", "Path to Things database (overrides THINGSDB)")
	cmd.Flags().StringVar(&dbPath, "database", "", "Alias for --db")
	addTaskQueryFlags(cmd, &opts, true, true)
	addTaskOutputFlags(cmd, &format, &selectRaw, &asJSON, &noHeader)

	return cmd
}
