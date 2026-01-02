package cli

import (
	"encoding/json"
	"fmt"

	"github.com/ossianhempel/things3-cli/internal/db"
	"github.com/spf13/cobra"
)

func NewAllCommand(app *App) *cobra.Command {
	var dbPath string
	var limit int
	var asJSON bool
	var noHeader bool
	var recursive bool

	cmd := &cobra.Command{
		Use:   "all",
		Short: "List key sections from the Things database",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, _, err := db.OpenDefault(dbPath)
			if err != nil {
				return formatDBError(err)
			}
			defer store.Close()

			incompleteFilter, _, err := buildTaskFilter(store, TaskQueryOptions{
				Status:           "incomplete",
				Limit:            limit,
				IncludeChecklist: recursive,
			})
			if err != nil {
				return err
			}
			anyFilter, _, err := buildTaskFilter(store, TaskQueryOptions{
				Status:           "any",
				Limit:            limit,
				IncludeChecklist: recursive,
			})
			if err != nil {
				return err
			}

			inbox, err := store.InboxTasks(incompleteFilter)
			if err != nil {
				return formatDBError(err)
			}
			today, err := store.TodayTasks(incompleteFilter)
			if err != nil {
				return formatDBError(err)
			}
			upcoming, err := store.UpcomingTasks(incompleteFilter)
			if err != nil {
				return formatDBError(err)
			}
			anytime, err := store.AnytimeTasks(incompleteFilter)
			if err != nil {
				return formatDBError(err)
			}
			someday, err := store.SomedayTasks(incompleteFilter)
			if err != nil {
				return formatDBError(err)
			}
			logbook, err := store.LogbookTasks(anyFilter)
			if err != nil {
				return formatDBError(err)
			}
			var noArea any
			var areas any
			if recursive {
				noAreaItems, err := store.ProjectsWithoutAreaTree(incompleteFilter, false)
				if err != nil {
					return formatDBError(err)
				}
				areaItems, err := store.AreasTree(incompleteFilter, false)
				if err != nil {
					return formatDBError(err)
				}
				noArea = noAreaItems
				areas = areaItems
			} else {
				noAreaItems, err := store.ProjectsWithoutArea(db.ProjectFilter{Status: incompleteFilter.Status, IncludeTrashed: incompleteFilter.IncludeTrashed})
				if err != nil {
					return formatDBError(err)
				}
				areaItems, err := store.Areas()
				if err != nil {
					return formatDBError(err)
				}
				noArea = noAreaItems
				areas = areaItems
			}

			if asJSON {
				sections := []struct {
					Title string `json:"title"`
					Items any    `json:"items"`
				}{
					{Title: "Inbox", Items: inbox},
					{Title: "Today", Items: today},
					{Title: "Upcoming", Items: upcoming},
					{Title: "Anytime", Items: anytime},
					{Title: "Someday", Items: someday},
					{Title: "Logbook", Items: logbook},
					{Title: "No Area", Items: noArea},
					{Title: "Areas", Items: areas},
				}
				enc := json.NewEncoder(app.Out)
				return enc.Encode(sections)
			}
			tableOpts := TaskOutputOptions{Format: "table", NoHeader: noHeader}

			first := true
			printSection := func(title string, fn func() error) error {
				if !first {
					fmt.Fprintln(app.Out)
				}
				first = false
				fmt.Fprintln(app.Out, title)
				return fn()
			}

			if err := printSection("Inbox", func() error {
				if len(inbox) == 0 {
					return nil
				}
				return printTasks(app.Out, inbox, tableOpts)
			}); err != nil {
				return err
			}
			if err := printSection("Today", func() error {
				if len(today) == 0 {
					return nil
				}
				return printTasks(app.Out, today, tableOpts)
			}); err != nil {
				return err
			}
			if err := printSection("Upcoming", func() error {
				if len(upcoming) == 0 {
					return nil
				}
				return printTasks(app.Out, upcoming, tableOpts)
			}); err != nil {
				return err
			}
			if err := printSection("Anytime", func() error {
				if len(anytime) == 0 {
					return nil
				}
				return printTasks(app.Out, anytime, tableOpts)
			}); err != nil {
				return err
			}
			if err := printSection("Someday", func() error {
				if len(someday) == 0 {
					return nil
				}
				return printTasks(app.Out, someday, tableOpts)
			}); err != nil {
				return err
			}
			if err := printSection("Logbook", func() error {
				if len(logbook) == 0 {
					return nil
				}
				return printTasks(app.Out, logbook, tableOpts)
			}); err != nil {
				return err
			}
			if err := printSection("No Area", func() error {
				if recursive {
					items := noArea.([]db.TreeItem)
					if len(items) == 0 {
						return nil
					}
					return printTree(app.Out, items, false)
				}
				items := noArea.([]db.Project)
				if len(items) == 0 {
					return nil
				}
				return printProjects(app.Out, items, false, noHeader)
			}); err != nil {
				return err
			}
			if err := printSection("Areas", func() error {
				if recursive {
					items := areas.([]db.TreeItem)
					if len(items) == 0 {
						return nil
					}
					return printTree(app.Out, items, false)
				}
				items := areas.([]db.Area)
				if len(items) == 0 {
					return nil
				}
				return printAreas(app.Out, items, false, noHeader)
			}); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&dbPath, "db", "d", "", "Path to Things database (overrides THINGSDB)")
	cmd.Flags().StringVar(&dbPath, "database", "", "Alias for --db")
	cmd.Flags().IntVar(&limit, "limit", 200, "Limit number of results (0 = no limit)")
	cmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Include checklist items in JSON output")
	cmd.Flags().BoolVarP(&asJSON, "json", "j", false, "Output JSON")
	cmd.Flags().BoolVar(&noHeader, "no-header", false, "Suppress header row")

	return cmd
}
