package cli

import (
	"github.com/ossianhempel/things3-cli/internal/db"
	"github.com/spf13/cobra"
)

func NewInboxCommand(app *App) *cobra.Command {
	return newTaskListCommand(app, "inbox", "List inbox tasks from the Things database", "incomplete", func(store *db.Store, filter db.TaskFilter) ([]db.Task, error) {
		return store.InboxTasks(filter)
	})
}

func NewAnytimeCommand(app *App) *cobra.Command {
	return newTaskListCommand(app, "anytime", "List Anytime tasks from the Things database", "incomplete", func(store *db.Store, filter db.TaskFilter) ([]db.Task, error) {
		return store.AnytimeTasks(filter)
	})
}

func NewSomedayCommand(app *App) *cobra.Command {
	return newTaskListCommand(app, "someday", "List Someday tasks from the Things database", "incomplete", func(store *db.Store, filter db.TaskFilter) ([]db.Task, error) {
		return store.SomedayTasks(filter)
	})
}

func NewUpcomingCommand(app *App) *cobra.Command {
	return newTaskListCommand(app, "upcoming", "List upcoming tasks from the Things database", "incomplete", func(store *db.Store, filter db.TaskFilter) ([]db.Task, error) {
		return store.UpcomingTasks(filter)
	})
}

func NewRepeatingCommand(app *App) *cobra.Command {
	var dbPath string
	opts := TaskQueryOptions{
		Status:        "incomplete",
		Limit:         200,
		RepeatingOnly: true,
	}
	var format string
	var selectRaw string
	var asJSON bool
	var noHeader bool

	cmd := &cobra.Command{
		Use:   "repeating",
		Short: "List repeating tasks and projects from the Things database",
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
			tasks, err := fetchTasks(store, store.Tasks, opts, false, []int{db.TaskTypeTodo, db.TaskTypeProject})
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

func NewTemplatesCommand(app *App) *cobra.Command {
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
		Use:   "templates",
		Short: "List repeating template tasks and projects from the Things database",
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
			if outputOpts.Format == "table" && len(outputOpts.Select) == 0 {
				outputOpts.Select = []string{"uuid", "title", "repeat_active", "repeat_mode", "repeat_unit", "repeat_interval", "repeat_anchor", "repeat_next_date"}
			}
			tasks, err := fetchTasks(store, store.TemplatesTasks, opts, false, []int{db.TaskTypeTodo, db.TaskTypeProject})
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

func NewDeadlinesCommand(app *App) *cobra.Command {
	return newTaskListCommand(app, "deadlines", "List tasks with deadlines from the Things database", "incomplete", func(store *db.Store, filter db.TaskFilter) ([]db.Task, error) {
		return store.DeadlinesTasks(filter)
	})
}

func NewLogbookCommand(app *App) *cobra.Command {
	return newTaskListCommand(app, "logbook", "List logbook tasks from the Things database", "any", func(store *db.Store, filter db.TaskFilter) ([]db.Task, error) {
		return store.LogbookTasks(filter)
	})
}

func NewCompletedCommand(app *App) *cobra.Command {
	return newTaskListCommand(app, "completed", "List completed tasks from the Things database", "completed", func(store *db.Store, filter db.TaskFilter) ([]db.Task, error) {
		return store.CompletedTasks(filter)
	})
}

func NewCanceledCommand(app *App) *cobra.Command {
	return newTaskListCommand(app, "canceled", "List canceled tasks from the Things database", "canceled", func(store *db.Store, filter db.TaskFilter) ([]db.Task, error) {
		return store.CanceledTasks(filter)
	})
}

func NewTrashCommand(app *App) *cobra.Command {
	return newTaskListCommand(app, "trash", "List trashed tasks from the Things database", "any", func(store *db.Store, filter db.TaskFilter) ([]db.Task, error) {
		return store.TrashTasks(filter)
	})
}
