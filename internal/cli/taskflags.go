package cli

import "github.com/spf13/cobra"

func addTaskQueryFlags(cmd *cobra.Command, opts *TaskQueryOptions, includeSearch bool, includeQuery bool) {
	flags := cmd.Flags()
	flags.StringVar(&opts.Status, "status", opts.Status, "Filter by status: incomplete, completed, canceled, any")
	flags.StringVarP(&opts.Project, "filter-project", "p", "", "Filter by project title or ID")
	flags.StringVar(&opts.Project, "project", "", "Alias for --filter-project")
	flags.StringVarP(&opts.Area, "filter-area", "a", "", "Filter by area title or ID")
	flags.StringVar(&opts.Area, "area", "", "Alias for --filter-area")
	flags.StringVarP(&opts.Tag, "filter-tag", "t", "", "Filter by tag title or ID")
	flags.StringVar(&opts.Tag, "filtertag", "", "Alias for --filter-tag")
	flags.StringVar(&opts.Tag, "tag", "", "Alias for --filter-tag")
	if includeSearch {
		flags.StringVar(&opts.Search, "search", "", "Search title or notes (case-insensitive substring)")
	}
	if includeQuery {
		flags.StringVar(&opts.Query, "query", "", "Rich query (boolean, fields, regex; e.g. title:/regex/ AND tag:reading)")
	}
	flags.IntVar(&opts.Limit, "limit", opts.Limit, "Limit number of results (0 = no limit)")
	flags.IntVar(&opts.Offset, "offset", 0, "Offset results for pagination")
	flags.BoolVar(&opts.IncludeTrashed, "include-trashed", false, "Include trashed tasks")
	flags.BoolVar(&opts.All, "all", false, "Include completed, canceled, and trashed tasks")
	flags.BoolVarP(&opts.IncludeChecklist, "recursive", "r", false, "Include checklist items in JSON output")
	flags.StringVar(&opts.CreatedBefore, "created-before", "", "Filter tasks created before (YYYY-MM-DD or RFC3339)")
	flags.StringVar(&opts.CreatedAfter, "created-after", "", "Filter tasks created after (YYYY-MM-DD or RFC3339)")
	flags.StringVar(&opts.ModifiedBefore, "modified-before", "", "Filter tasks modified before (YYYY-MM-DD or RFC3339)")
	flags.StringVar(&opts.ModifiedAfter, "modified-after", "", "Filter tasks modified after (YYYY-MM-DD or RFC3339)")
	flags.StringVar(&opts.DueBefore, "due-before", "", "Filter tasks due before (YYYY-MM-DD)")
	flags.StringVar(&opts.StartBefore, "start-before", "", "Filter tasks starting before (YYYY-MM-DD)")
	flags.BoolVar(&opts.HasURL, "has-url", false, "Filter tasks with URLs in notes")
	flags.StringVar(&opts.Sort, "sort", "", "Sort by fields (e.g. created,-deadline,title)")
}

func addTaskOutputFlags(cmd *cobra.Command, format *string, selectRaw *string, asJSON *bool, noHeader *bool) {
	flags := cmd.Flags()
	flags.StringVar(format, "format", "", "Output format: table, json, jsonl, csv")
	flags.StringVar(selectRaw, "select", "", "Select fields (comma-separated)")
	flags.BoolVarP(asJSON, "json", "j", false, "Output JSON")
	flags.BoolVar(noHeader, "no-header", false, "Suppress header row")
}
