package cli

import (
	"fmt"
	"strings"

	"github.com/ossianhempel/things3-cli/internal/db"
	"github.com/ossianhempel/things3-cli/internal/things"
	"github.com/spf13/cobra"
)

// NewUpdateProjectCommand builds the update-project subcommand.
func NewUpdateProjectCommand(app *App) *cobra.Command {
	opts := things.UpdateProjectOptions{}
	repeatOpts := RepeatOptions{}
	var allowUnsafeTitle bool
	var dbPath string
	var repeatJSON bool

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
			if err := validateWhenInput(opts.When); err != nil {
				return err
			}

			repeatSpec, err := parseRepeatSpec(cmd, repeatOpts)
			if err != nil {
				return err
			}
			if repeatSpec.Enabled && strings.TrimSpace(opts.ID) == "" {
				return fmt.Errorf("Error: repeating updates require --id")
			}
			if repeatSpec.Enabled && (opts.Completed || opts.Canceled) {
				return fmt.Errorf("Error: --completed/--canceled cannot be combined with repeat changes because repeat templates must remain incomplete")
			}
			if repeatSpec.Enabled && opts.Duplicate {
				return fmt.Errorf("Error: --duplicate cannot be combined with repeat changes because the repeat rule would apply to the source template, not the duplicate")
			}
			if repeatJSON && !repeatSpec.Enabled {
				return fmt.Errorf("Error: --json is currently supported only for repeat updates")
			}
			if repeatSpec.Enabled && !repeatSpec.Clear && opts.Deadline != "" && repeatSpec.Spec.DeadlineOffset != nil {
				return fmt.Errorf("Error: --deadline cannot be combined with --repeat-deadline")
			}
			if repeatSpec.Enabled && !repeatSpec.Clear && (cmd.Flags().Changed("when")) {
				return fmt.Errorf("Error: --when/--later cannot be combined with --repeat because repeat activation controls scheduling; apply the repeat first, then update the template separately")
			}
			opts.NotesSet = cmd.Flags().Changed("notes")
			opts.WhenSet = cmd.Flags().Changed("when")
			opts.DeadlineSet = cmd.Flags().Changed("deadline")
			opts.TagsSet = cmd.Flags().Changed("tags")

			if !repeatSpec.Enabled {
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
			}

			hasChanges := hasProjectUpdateChanges(opts, rawInput)
			plan, err := prepareRepeatUpdateForType(dbPath, opts.ID, repeatSpec, app.DryRun, db.TaskTypeProject)
			if err != nil {
				return formatDBError(err)
			}
			defer plan.store.Close()
			if hasChanges && plan.usedTemplate {
				return fmt.Errorf("Error: combined ordinary and repeat updates cannot target a generated occurrence; rerun against template UUID %s", plan.targetID)
			}

			if hasChanges {
				if err := ensureProjectAuth(app, &opts); err != nil {
					return err
				}
				url, err := things.BuildUpdateProjectURL(opts, rawInput)
				if err != nil {
					return err
				}
				if app.DryRun {
					if !repeatJSON {
						if err := openURL(app, url); err != nil {
							return err
						}
					}
					plan.result.Intent = &repeatIntent{URL: redactURLIntent(url)}
					plan.result.addStage(repeatStageURL, repeatStatusPlanned)
					plan.result.addStage(repeatStageDatabase, repeatStatusPlanned)
					plan.result.addStage(repeatStageVerification, repeatStatusPlanned)
					return renderRepeatResult(app.Out, plan.result, repeatJSON)
				}
				if err := openURL(app, url); err != nil {
					return err
				}
				plan.result.Intent = &repeatIntent{URL: redactURLIntent(url)}
				plan.result.addStage(repeatStageURL, repeatStatusCompleted)
			} else if app.DryRun {
				plan.result.addStage(repeatStageDatabase, repeatStatusPlanned)
				plan.result.addStage(repeatStageVerification, repeatStatusPlanned)
				return renderRepeatResult(app.Out, plan.result, repeatJSON)
			}
			if plan.usedTemplate {
				fmt.Fprintf(app.Err, "Note: resolved repeating template %s for update\n", plan.targetID)
			}
			if err := plan.executeMutation(); err != nil {
				_ = renderRepeatResult(app.Out, plan.result, repeatJSON)
				return formatDBError(err)
			}
			return renderRepeatResult(app.Out, plan.result, repeatJSON)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&dbPath, "db", "d", "", "Path to Things database (overrides THINGSDB)")
	flags.StringVar(&dbPath, "database", "", "Alias for --db")
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
	flags.BoolVar(&repeatJSON, "json", false, "Emit a structured result for repeat updates")
	addRepeatFlags(cmd, &repeatOpts, true)

	return cmd
}

func ensureProjectAuth(app *App, opts *things.UpdateProjectOptions) error {
	token, err := resolveAuthToken(app, opts.AuthToken)
	if err != nil {
		return err
	}
	opts.AuthToken = token
	return nil
}

func hasProjectUpdateChanges(opts things.UpdateProjectOptions, rawInput string) bool {
	if strings.TrimSpace(rawInput) != "" {
		return true
	}
	if opts.Notes != "" || opts.NotesSet || opts.PrependNotes != "" || opts.AppendNotes != "" {
		return true
	}
	if opts.When != "" || opts.WhenSet {
		return true
	}
	if opts.Deadline != "" || opts.DeadlineSet {
		return true
	}
	if opts.Tags != "" || opts.TagsSet || opts.AddTags != "" {
		return true
	}
	if opts.Completed || opts.Canceled {
		return true
	}
	if opts.Reveal || opts.Duplicate {
		return true
	}
	if opts.CompletionDate != "" || opts.CreationDate != "" {
		return true
	}
	if opts.AreaID != "" || opts.Area != "" {
		return true
	}
	if len(opts.Todos) > 0 {
		return true
	}
	return false
}
