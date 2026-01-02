package cli

import (
	"time"

	"github.com/ossianhempel/things3-cli/internal/db"
	"github.com/spf13/cobra"
)

func NewCreatedTodayCommand(app *App) *cobra.Command {
	var dbPath string
	opts := TaskQueryOptions{
		Status: "any",
		Limit:  200,
	}
	var format string
	var selectRaw string
	var asJSON bool
	var noHeader bool

	cmd := &cobra.Command{
		Use:   "createdtoday",
		Short: "List tasks created today from the Things database",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, _, err := db.OpenDefault(dbPath)
			if err != nil {
				return formatDBError(err)
			}
			defer store.Close()

			now := time.Now()
			start := now.Add(-24 * time.Hour)
			opts.HasURLSet = cmd.Flags().Changed("has-url")
			outputOpts, err := resolveTaskOutputOptions(format, asJSON, selectRaw, noHeader)
			if err != nil {
				return err
			}
			forcePost := opts.Query != "" || opts.Sort != "" || opts.Offset > 0
			tasks, err := fetchTasks(store, func(filter db.TaskFilter) ([]db.Task, error) {
				return store.TasksCreatedBetween(start, now, filter)
			}, opts, forcePost, []int{db.TaskTypeTodo})
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

func NewLogTodayCommand(app *App) *cobra.Command {
	var dbPath string
	opts := TaskQueryOptions{
		Status: "any",
		Limit:  200,
	}
	var format string
	var selectRaw string
	var asJSON bool
	var noHeader bool

	cmd := &cobra.Command{
		Use:   "logtoday",
		Short: "List tasks completed today from the Things database",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, _, err := db.OpenDefault(dbPath)
			if err != nil {
				return formatDBError(err)
			}
			defer store.Close()

			start, end := dayBounds()
			opts.HasURLSet = cmd.Flags().Changed("has-url")
			outputOpts, err := resolveTaskOutputOptions(format, asJSON, selectRaw, noHeader)
			if err != nil {
				return err
			}
			forcePost := opts.Query != "" || opts.Sort != "" || opts.Offset > 0
			tasks, err := fetchTasks(store, func(filter db.TaskFilter) ([]db.Task, error) {
				return store.TasksCompletedBetween(start, end, filter)
			}, opts, forcePost, []int{db.TaskTypeTodo})
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

func dayBounds() (time.Time, time.Time) {
	now := time.Now()
	loc := now.Location()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	end := start.Add(24 * time.Hour)
	return start, end
}
