package cli

import (
	"github.com/ossianhempel/things3-cli/internal/things"
	"github.com/spf13/cobra"
)

// NewAddProjectCommand builds the add-project subcommand.
func NewAddProjectCommand(app *App) *cobra.Command {
	opts := things.AddProjectOptions{}
	var allowUnsafeTitle bool

	cmd := &cobra.Command{
		Use:     "add-project [OPTIONS...] [-|TITLE]",
		Aliases: []string{"create-project"},
		Short:   "Add a new project",
		RunE: func(cmd *cobra.Command, args []string) error {
			rawInput, err := readInput(app.In, args)
			if err != nil {
				return err
			}
			title := extractTitle(rawInput, "")
			if err := guardUnsafeTitle(title, allowUnsafeTitle); err != nil {
				return err
			}
			if err := validateWhenInput(opts.When); err != nil {
				return err
			}

			url := things.BuildAddProjectURL(opts, rawInput)
			return openURL(app, url)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.AreaID, "area-id", "", "Area ID to add to")
	flags.StringVar(&opts.Area, "area", "", "Area to add to")
	flags.BoolVar(&opts.Canceled, "canceled", false, "Mark the project canceled")
	flags.BoolVar(&opts.Canceled, "cancelled", false, "Mark the project cancelled")
	flags.BoolVar(&opts.Completed, "completed", false, "Mark the project completed")
	flags.StringVar(&opts.CompletionDate, "completion-date", "", "Completion date (ISO8601)")
	flags.StringVar(&opts.CreationDate, "creation-date", "", "Creation date (ISO8601)")
	flags.StringVar(&opts.Deadline, "deadline", "", "Deadline for the project")
	flags.StringVar(&opts.Notes, "notes", "", "Notes for the project")
	flags.BoolVar(&opts.Reveal, "reveal", false, "Reveal the newly created project")
	flags.StringVar(&opts.Tags, "tags", "", "Comma-separated tags")
	flags.StringVar(&opts.When, "when", "", "When to schedule the project")
	flags.StringArrayVar(&opts.Todos, "todo", nil, "Todo title to add (repeatable)")
	flags.BoolVar(&allowUnsafeTitle, "allow-unsafe-title", false, "Allow titles that look like flag assignments")

	return cmd
}
