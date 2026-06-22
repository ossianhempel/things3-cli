package cli

import (
	"time"

	"github.com/ossianhempel/things3-cli/internal/things"
	"github.com/spf13/cobra"
)

// NewUpdateProjectCommand builds the update-project subcommand.
func NewUpdateProjectCommand(app *App) *cobra.Command {
	opts := things.UpdateProjectOptions{}
	var allowUnsafeTitle bool

	cmd := &cobra.Command{
		Use:   "update-project [OPTIONS...] [--] [-|TITLE]",
		Short: "Update an existing project",
		RunE: func(cmd *cobra.Command, args []string) error {
			rawInput, err := readInput(app.In, args)
			if err != nil {
				return err
			}
			title := extractTitle(rawInput, "")
			if err := guardUnsafeTitle(title, allowUnsafeTitle); err != nil {
				return err
			}
			opts.NotesSet = cmd.Flags().Changed("notes")
			opts.WhenSet = cmd.Flags().Changed("when")
			opts.DeadlineSet = cmd.Flags().Changed("deadline")
			opts.TagsSet = cmd.Flags().Changed("tags")
			now := time.Now()
			opts.When, err = normalizeWhenInput(opts.When, now)
			if err != nil {
				return err
			}
			opts.Deadline, err = normalizeDeadlineInput(opts.Deadline, now)
			if err != nil {
				return err
			}

			token, err := resolveAuthToken(app, opts.AuthToken)
			if err != nil {
				return err
			}
			opts.AuthToken = token

			url, err := things.BuildUpdateProjectURL(opts, rawInput)
			if err != nil {
				return err
			}
			return openURL(app, url)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.AuthToken, "auth-token", "", "Things URL scheme authorization token")
	flags.StringVar(&opts.ID, "id", "", "ID of the project to update")
	flags.StringVar(&opts.Notes, "notes", "", "Replace notes")
	flags.StringVar(&opts.PrependNotes, "prepend-notes", "", "Prepend to notes")
	flags.StringVar(&opts.AppendNotes, "append-notes", "", "Append to notes")
	flags.StringVar(&opts.When, "when", "", "When to schedule the project (empty string clears the date)")
	flags.StringVar(&opts.Deadline, "deadline", "", "Deadline for the project (empty string clears it)")
	flags.StringVar(&opts.Tags, "tags", "", "Replace tags (empty string clears all tags)")
	flags.StringVar(&opts.AddTags, "add-tags", "", "Add tags")
	flags.StringVar(&opts.AreaID, "area-id", "", "Area ID to move to")
	flags.StringVar(&opts.Area, "area", "", "Area to move to")
	flags.BoolVar(&opts.Completed, "completed", false, "Mark the project completed")
	flags.BoolVar(&opts.Canceled, "canceled", false, "Mark the project canceled")
	flags.BoolVar(&opts.Canceled, "cancelled", false, "Mark the project cancelled")
	flags.BoolVar(&opts.Reveal, "reveal", false, "Reveal the updated project")
	flags.BoolVar(&opts.Duplicate, "duplicate", false, "Duplicate before updating")
	flags.StringVar(&opts.CompletionDate, "completion-date", "", "Completion date (ISO8601)")
	flags.StringVar(&opts.CreationDate, "creation-date", "", "Creation date (ISO8601)")
	flags.StringArrayVar(&opts.Todos, "todo", nil, "Todo title to add (repeatable)")
	flags.BoolVar(&allowUnsafeTitle, "allow-unsafe-title", false, "Allow titles that look like flag assignments")

	return cmd
}
