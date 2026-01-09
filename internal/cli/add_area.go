package cli

import (
	"github.com/ossianhempel/things3-cli/internal/things"
	"github.com/spf13/cobra"
)

// NewAddAreaCommand builds the add-area subcommand.
func NewAddAreaCommand(app *App) *cobra.Command {
	opts := things.AddAreaOptions{}
	var allowUnsafeTitle bool

	cmd := &cobra.Command{
		Use:     "add-area [OPTIONS...] [-|TITLE]",
		Aliases: []string{"create-area"},
		Short:   "Add a new area",
		RunE: func(cmd *cobra.Command, args []string) error {
			rawInput, err := readInput(app.In, args)
			if err != nil {
				return err
			}
			title := extractTitle(rawInput, "")
			if err := guardUnsafeTitle(title, allowUnsafeTitle); err != nil {
				return err
			}

			script, err := things.BuildAddAreaScript(opts, rawInput)
			if err != nil {
				return err
			}
			return runScript(app, script)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.Tags, "tags", "", "Comma-separated tags")
	flags.BoolVar(&allowUnsafeTitle, "allow-unsafe-title", false, "Allow titles that look like flag assignments")

	return cmd
}
