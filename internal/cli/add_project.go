package cli

import (
	"fmt"
	"time"

	"github.com/ossianhempel/things3-cli/internal/db"
	"github.com/ossianhempel/things3-cli/internal/repeat"
	"github.com/ossianhempel/things3-cli/internal/things"
	"github.com/spf13/cobra"
)

// NewAddProjectCommand builds the add-project subcommand.
func NewAddProjectCommand(app *App) *cobra.Command {
	opts := things.AddProjectOptions{}
	repeatOpts := RepeatOptions{}
	var allowUnsafeTitle bool
	var dbPath string
	var repeatJSON bool

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

			repeatSpec, err := parseRepeatSpec(cmd, repeatOpts)
			if err != nil {
				return err
			}
			if repeatJSON && !repeatSpec.Enabled {
				return fmt.Errorf("Error: --json is currently supported only for repeat adds")
			}
			if repeatSpec.Enabled {
				if repeatSpec.Clear {
					return fmt.Errorf("Error: --repeat-clear is only valid with update commands")
				}
				if title == "" {
					return fmt.Errorf("Error: repeating add requires an explicit title")
				}
				if opts.Canceled || opts.Completed {
					return fmt.Errorf("Error: --completed/--canceled cannot be combined with repeat changes because repeat templates must remain incomplete")
				}
				if opts.Deadline != "" && repeatSpec.Spec.DeadlineOffset != nil {
					return fmt.Errorf("Error: --deadline cannot be combined with --repeat-deadline")
				}
				if cmd.Flags().Changed("when") {
					return fmt.Errorf("Error: --when cannot be combined with --repeat because repeat activation controls scheduling; create the repeat first, then inspect its template state")
				}
			}

			url := things.BuildAddProjectURL(opts, rawInput)
			if !repeatSpec.Enabled {
				return openURL(app, url)
			}
			update, err := repeat.BuildUpdate(repeatSpec.Spec)
			if err != nil {
				return err
			}
			if app.DryRun {
				store, resolvedPath, err := db.OpenDefault(dbPath)
				if err != nil {
					return formatDBError(err)
				}
				defer store.Close()
				if err := store.ValidateRepeatSchema(); err != nil {
					return formatDBError(err)
				}
				result := repeatResult{SchemaVersion: 1, Action: "create", DryRun: true, Repeat: expectedRepeatState(repeatSpec.Spec), Database: repeatDatabase{Path: store.Path(), Source: repeatDatabaseSource(dbPath, resolvedPath)}, Intent: &repeatIntent{URL: redactURLIntent(url)}}
				result.addStage(repeatStageURL, repeatStatusPlanned)
				result.addStage(repeatStageLocate, repeatStatusPlanned)
				result.addStage(repeatStageDatabase, repeatStatusPlanned)
				result.addStage(repeatStageVerification, repeatStatusPlanned)
				return renderRepeatResult(app.Out, result, repeatJSON)
			}

			expected := expectedRepeatState(repeatSpec.Spec)
			result := repeatResult{SchemaVersion: 1, Action: "create", Repeat: expected}
			store, resolvedPath, err := db.OpenDefaultWritable(dbPath)
			if err != nil {
				return formatDBError(err)
			}
			defer store.Close()
			result.Database = repeatDatabase{Path: store.Path(), Source: repeatDatabaseSource(dbPath, resolvedPath)}
			result.Intent = &repeatIntent{URL: redactURLIntent(url)}
			if err := store.ValidateRepeatSchema(); err != nil {
				return formatDBError(err)
			}
			ensureThingsLaunched(app)
			started := time.Now().Add(-2 * time.Second)
			if err := openURL(app, url); err != nil {
				return err
			}
			result.addStage(repeatStageURL, repeatStatusCompleted)

			projectID, err := waitForCreatedItem(store, title, db.TaskTypeProject, started)
			if err != nil {
				result.failStage(repeatStageLocate)
				result.Recovery = []repeatRecovery{{Argv: []string{"things", "search", title, "--select", "uuid,title"}}}
				_ = renderRepeatResult(app.Out, result, repeatJSON)
				return formatDBError(err)
			}
			result.IDs.Created = projectID
			result.IDs.Template = projectID
			result.addStage(repeatStageLocate, repeatStatusCompleted)
			if err := store.ApplyRepeatRule(projectID, update); err != nil {
				result.failStage(repeatStageDatabase)
				result.Recovery = []repeatRecovery{{Argv: recoveryArgvForType(projectID, repeatSpec, result.Database.Path, db.TaskTypeProject)}}
				_ = renderRepeatResult(app.Out, result, repeatJSON)
				return formatDBError(err)
			}
			result.addStage(repeatStageDatabase, repeatStatusApplied)
			actual, err := verifyRepeatState(store, projectID, expected, false)
			result.Repeat = actual
			if err != nil {
				result.failStage(repeatStageVerification)
				result.Recovery = []repeatRecovery{{Argv: recoveryArgvForType(projectID, repeatSpec, result.Database.Path, db.TaskTypeProject)}}
				_ = renderRepeatResult(app.Out, result, repeatJSON)
				return formatDBError(err)
			}
			result.markVerified()
			return renderRepeatResult(app.Out, result, repeatJSON)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&dbPath, "db", "d", "", "Path to Things database (overrides THINGSDB)")
	flags.StringVar(&dbPath, "database", "", "Alias for --db")
	flags.BoolVar(&repeatJSON, "json", false, "Emit a structured result for repeat adds")
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
	addRepeatFlags(cmd, &repeatOpts, false)

	return cmd
}
