package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewRoot builds the root command.
func NewRoot(app *App) *cobra.Command {
	var versionFlag bool

	cmd := &cobra.Command{
		Use:           "things",
		Short:         "Manage Things 3 from the terminal",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if versionFlag {
				printVersion(app.Out)
				return ErrVersionPrinted
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			printHelp(app.Out, formatHelpText(rootHelp, isTTY(app.Out)))
			return ErrHelpPrinted
		},
	}

	cmd.PersistentFlags().BoolVar(&app.Debug, "debug", false, "Enable debug mode")
	cmd.PersistentFlags().BoolVar(&app.Foreground, "foreground", false, "Open Things in the foreground")
	cmd.PersistentFlags().BoolVar(&app.DryRun, "dry-run", false, "Print the Things URL without opening it")
	cmd.PersistentFlags().BoolVarP(&versionFlag, "version", "V", false, "Print version information")

	cmd.AddCommand(NewAddCommand(app))
	cmd.AddCommand(NewAddAreaCommand(app))
	cmd.AddCommand(NewAddProjectCommand(app))
	cmd.AddCommand(NewAreasCommand(app))
	cmd.AddCommand(NewInboxCommand(app))
	cmd.AddCommand(NewTodayCommand(app))
	cmd.AddCommand(NewUpcomingCommand(app))
	cmd.AddCommand(NewAnytimeCommand(app))
	cmd.AddCommand(NewSomedayCommand(app))
	cmd.AddCommand(NewLogbookCommand(app))
	cmd.AddCommand(NewLogTodayCommand(app))
	cmd.AddCommand(NewCreatedTodayCommand(app))
	cmd.AddCommand(NewCompletedCommand(app))
	cmd.AddCommand(NewCanceledCommand(app))
	cmd.AddCommand(NewTrashCommand(app))
	cmd.AddCommand(NewDeadlinesCommand(app))
	cmd.AddCommand(NewAllCommand(app))
	cmd.AddCommand(NewProjectsCommand(app))
	cmd.AddCommand(NewTagsCommand(app))
	cmd.AddCommand(NewTasksCommand(app))
	cmd.AddCommand(NewAuthCommand(app))
	cmd.AddCommand(NewUpdateCommand(app))
	cmd.AddCommand(NewUpdateAreaCommand(app))
	cmd.AddCommand(NewUpdateProjectCommand(app))
	cmd.AddCommand(NewDeleteCommand(app))
	cmd.AddCommand(NewUndoCommand(app))
	cmd.AddCommand(NewShowCommand(app))
	cmd.AddCommand(NewSearchCommand(app))

	cmd.SetHelpCommand(&cobra.Command{
		Use:   "help [command]",
		Short: "Show help for things and its subcommands",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				printHelp(app.Out, formatHelpText(rootHelp, isTTY(app.Out)))
				return ErrHelpPrinted
			}
			switch args[0] {
			case "add":
				printHelp(app.Out, formatHelpText(addHelp, isTTY(app.Out)))
			case "add-area", "create-area":
				printHelp(app.Out, formatHelpText(addAreaHelp, isTTY(app.Out)))
			case "add-project":
				printHelp(app.Out, formatHelpText(addProjectHelp, isTTY(app.Out)))
			case "create-project":
				printHelp(app.Out, formatHelpText(addProjectHelp, isTTY(app.Out)))
			case "areas":
				printHelp(app.Out, formatHelpText(areasHelp, isTTY(app.Out)))
			case "inbox":
				printHelp(app.Out, formatHelpText(inboxHelp, isTTY(app.Out)))
			case "today":
				printHelp(app.Out, formatHelpText(todayHelp, isTTY(app.Out)))
			case "upcoming":
				printHelp(app.Out, formatHelpText(upcomingHelp, isTTY(app.Out)))
			case "anytime":
				printHelp(app.Out, formatHelpText(anytimeHelp, isTTY(app.Out)))
			case "someday":
				printHelp(app.Out, formatHelpText(somedayHelp, isTTY(app.Out)))
			case "logbook":
				printHelp(app.Out, formatHelpText(logbookHelp, isTTY(app.Out)))
			case "logtoday":
				printHelp(app.Out, formatHelpText(logtodayHelp, isTTY(app.Out)))
			case "createdtoday":
				printHelp(app.Out, formatHelpText(createdTodayHelp, isTTY(app.Out)))
			case "completed":
				printHelp(app.Out, formatHelpText(completedHelp, isTTY(app.Out)))
			case "canceled":
				printHelp(app.Out, formatHelpText(canceledHelp, isTTY(app.Out)))
			case "trash":
				printHelp(app.Out, formatHelpText(trashHelp, isTTY(app.Out)))
			case "deadlines":
				printHelp(app.Out, formatHelpText(deadlinesHelp, isTTY(app.Out)))
			case "all":
				printHelp(app.Out, formatHelpText(allHelp, isTTY(app.Out)))
			case "projects":
				printHelp(app.Out, formatHelpText(projectsHelp, isTTY(app.Out)))
			case "tags":
				printHelp(app.Out, formatHelpText(tagsHelp, isTTY(app.Out)))
			case "tasks":
				printHelp(app.Out, formatHelpText(tasksHelp, isTTY(app.Out)))
			case "auth":
				printHelp(app.Out, formatHelpText(authHelp, isTTY(app.Out)))
			case "show":
				printHelp(app.Out, formatHelpText(showHelp, isTTY(app.Out)))
			case "search":
				printHelp(app.Out, formatHelpText(searchHelp, isTTY(app.Out)))
			case "update":
				printHelp(app.Out, formatHelpText(updateHelp, isTTY(app.Out)))
			case "update-area":
				printHelp(app.Out, formatHelpText(updateAreaHelp, isTTY(app.Out)))
			case "update-project":
				printHelp(app.Out, formatHelpText(updateProjectHelp, isTTY(app.Out)))
			case "delete":
				printHelp(app.Out, formatHelpText(deleteHelp, isTTY(app.Out)))
			case "undo":
				printHelp(app.Out, formatHelpText(undoHelp, isTTY(app.Out)))
			case "help":
				printHelp(app.Out, formatHelpText(rootHelp, isTTY(app.Out)))
			default:
				return fmt.Errorf("Error: Invalid command `things %s'", args[0])
			}
			return ErrHelpPrinted
		},
	})

	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		name := cmd.Name()
		switch name {
		case "things":
			printHelp(app.Out, formatHelpText(rootHelp, isTTY(app.Out)))
		case "add":
			printHelp(app.Out, formatHelpText(addHelp, isTTY(app.Out)))
		case "add-area":
			printHelp(app.Out, formatHelpText(addAreaHelp, isTTY(app.Out)))
		case "add-project":
			printHelp(app.Out, formatHelpText(addProjectHelp, isTTY(app.Out)))
		case "create-project":
			printHelp(app.Out, formatHelpText(addProjectHelp, isTTY(app.Out)))
		case "create-area":
			printHelp(app.Out, formatHelpText(addAreaHelp, isTTY(app.Out)))
		case "areas":
			printHelp(app.Out, formatHelpText(areasHelp, isTTY(app.Out)))
		case "inbox":
			printHelp(app.Out, formatHelpText(inboxHelp, isTTY(app.Out)))
		case "today":
			printHelp(app.Out, formatHelpText(todayHelp, isTTY(app.Out)))
		case "upcoming":
			printHelp(app.Out, formatHelpText(upcomingHelp, isTTY(app.Out)))
		case "anytime":
			printHelp(app.Out, formatHelpText(anytimeHelp, isTTY(app.Out)))
		case "someday":
			printHelp(app.Out, formatHelpText(somedayHelp, isTTY(app.Out)))
		case "logbook":
			printHelp(app.Out, formatHelpText(logbookHelp, isTTY(app.Out)))
		case "logtoday":
			printHelp(app.Out, formatHelpText(logtodayHelp, isTTY(app.Out)))
		case "createdtoday":
			printHelp(app.Out, formatHelpText(createdTodayHelp, isTTY(app.Out)))
		case "completed":
			printHelp(app.Out, formatHelpText(completedHelp, isTTY(app.Out)))
		case "canceled":
			printHelp(app.Out, formatHelpText(canceledHelp, isTTY(app.Out)))
		case "trash":
			printHelp(app.Out, formatHelpText(trashHelp, isTTY(app.Out)))
		case "deadlines":
			printHelp(app.Out, formatHelpText(deadlinesHelp, isTTY(app.Out)))
		case "all":
			printHelp(app.Out, formatHelpText(allHelp, isTTY(app.Out)))
		case "projects":
			printHelp(app.Out, formatHelpText(projectsHelp, isTTY(app.Out)))
		case "tags":
			printHelp(app.Out, formatHelpText(tagsHelp, isTTY(app.Out)))
		case "tasks":
			printHelp(app.Out, formatHelpText(tasksHelp, isTTY(app.Out)))
		case "auth":
			printHelp(app.Out, formatHelpText(authHelp, isTTY(app.Out)))
		case "show":
			printHelp(app.Out, formatHelpText(showHelp, isTTY(app.Out)))
		case "search":
			printHelp(app.Out, formatHelpText(searchHelp, isTTY(app.Out)))
		case "update":
			printHelp(app.Out, formatHelpText(updateHelp, isTTY(app.Out)))
		case "update-area":
			printHelp(app.Out, formatHelpText(updateAreaHelp, isTTY(app.Out)))
		case "update-project":
			printHelp(app.Out, formatHelpText(updateProjectHelp, isTTY(app.Out)))
		case "delete":
			printHelp(app.Out, formatHelpText(deleteHelp, isTTY(app.Out)))
		case "undo":
			printHelp(app.Out, formatHelpText(undoHelp, isTTY(app.Out)))
		default:
			printHelp(app.Out, formatHelpText(rootHelp, isTTY(app.Out)))
		}
	})

	return cmd
}
