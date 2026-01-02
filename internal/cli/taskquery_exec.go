package cli

import "github.com/ossianhempel/things3-cli/internal/db"

func fetchTasks(store *db.Store, runner func(db.TaskFilter) ([]db.Task, error), opts TaskQueryOptions, forcePost bool, types []int) ([]db.Task, error) {
	filter, sortSpec, err := buildTaskFilter(store, opts)
	if err != nil {
		return nil, err
	}
	if len(types) > 0 {
		filter.Types = types
	}

	queryExpr, err := parseRichQuery(opts.Query)
	if err != nil {
		return nil, err
	}

	postProcess := forcePost || queryExpr != nil
	if postProcess {
		filter.Limit = 0
		filter.Offset = 0
	}

	tasks, err := runner(filter)
	if err != nil {
		return nil, err
	}

	if queryExpr != nil {
		tasks = filterTasksByQuery(tasks, queryExpr)
	}

	if postProcess && len(sortSpec) > 0 {
		sortTasks(tasks, sortSpec)
	}

	if postProcess {
		tasks = applyOffsetLimit(tasks, opts.Limit, opts.Offset)
	}

	return tasks, nil
}

func applyOffsetLimit(tasks []db.Task, limit int, offset int) []db.Task {
	if offset < 0 {
		offset = 0
	}
	if offset >= len(tasks) {
		return []db.Task{}
	}
	if offset > 0 {
		tasks = tasks[offset:]
	}
	if limit > 0 && len(tasks) > limit {
		return tasks[:limit]
	}
	return tasks
}

func hasExplicitSelector(cmdFlags map[string]bool, opts TaskQueryOptions) bool {
	if opts.Search != "" || opts.Query != "" || opts.Project != "" || opts.Area != "" || opts.Tag != "" {
		return true
	}
	if opts.CreatedBefore != "" || opts.CreatedAfter != "" || opts.ModifiedBefore != "" || opts.ModifiedAfter != "" {
		return true
	}
	if opts.DueBefore != "" || opts.StartBefore != "" {
		return true
	}
	if opts.HasURLSet {
		return true
	}
	if opts.All || opts.IncludeTrashed {
		return true
	}
	if cmdFlags != nil && cmdFlags["status"] {
		return true
	}
	return false
}
