package cli

import (
	"fmt"
	"strings"

	"github.com/ossianhempel/things3-cli/internal/db"
	"github.com/spf13/cobra"
)

// NewSearchCommand builds the search subcommand.
func NewSearchCommand(app *App) *cobra.Command {
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
		Use:   "search [--] [-|QUERY]",
		Short: "Search tasks in the Things database",
		RunE: func(cmd *cobra.Command, args []string) error {
			query, err := readInput(app.In, args)
			if err != nil {
				return err
			}
			query = strings.TrimSpace(query)
			if query == "" && strings.TrimSpace(opts.Query) == "" {
				return fmt.Errorf("Error: query required")
			}
			if query != "" && strings.TrimSpace(opts.Query) != "" {
				return fmt.Errorf("Error: use either QUERY argument or --query")
			}
			if query != "" {
				opts.Search = query
			}

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
	addTaskQueryFlags(cmd, &opts, false, true)
	addTaskOutputFlags(cmd, &format, &selectRaw, &asJSON, &noHeader)

	return cmd
}
