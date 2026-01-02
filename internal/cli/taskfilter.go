package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/ossianhempel/things3-cli/internal/db"
)

type TaskQueryOptions struct {
	Status           string
	IncludeTrashed   bool
	All              bool
	Project          string
	Area             string
	Tag              string
	Search           string
	Query            string
	Limit            int
	Offset           int
	IncludeChecklist bool
	CreatedBefore    string
	CreatedAfter     string
	ModifiedBefore   string
	ModifiedAfter    string
	DueBefore        string
	StartBefore      string
	HasURL           bool
	HasURLSet        bool
	Sort             string
}

type TaskSortField struct {
	Field string
	Desc  bool
}

func buildTaskFilter(store *db.Store, opts TaskQueryOptions) (db.TaskFilter, []TaskSortField, error) {
	statusFilter, err := db.ParseStatus(opts.Status)
	if err != nil {
		return db.TaskFilter{}, nil, fmt.Errorf("Error: %s", err)
	}
	includeTrashed := opts.IncludeTrashed
	if opts.All {
		statusFilter = nil
		includeTrashed = true
	}

	projectID := ""
	if opts.Project != "" {
		projectID, err = store.ResolveProjectID(opts.Project)
		if err != nil {
			return db.TaskFilter{}, nil, fmt.Errorf("Error: %s", err)
		}
	}

	areaID := ""
	if opts.Area != "" {
		areaID, err = store.ResolveAreaID(opts.Area)
		if err != nil {
			return db.TaskFilter{}, nil, fmt.Errorf("Error: %s", err)
		}
	}

	tagID := ""
	if opts.Tag != "" {
		tagID, err = store.ResolveTagID(opts.Tag)
		if err != nil {
			return db.TaskFilter{}, nil, fmt.Errorf("Error: %s", err)
		}
	}

	sortFields, orderClause, err := parseSortSpec(opts.Sort)
	if err != nil {
		return db.TaskFilter{}, nil, err
	}

	filter := db.TaskFilter{
		Status:                statusFilter,
		IncludeTrashed:        includeTrashed,
		ExcludeTrashedContext: true,
		ProjectID:             projectID,
		AreaID:                areaID,
		TagID:                 tagID,
		Search:                opts.Search,
		Limit:                 opts.Limit,
		Offset:                opts.Offset,
		IncludeChecklist:      opts.IncludeChecklist,
		Order:                 orderClause,
	}

	if opts.CreatedAfter != "" {
		value, err := parseTimestampBound(opts.CreatedAfter, false)
		if err != nil {
			return db.TaskFilter{}, nil, err
		}
		filter.CreatedAfter = &value
	}
	if opts.CreatedBefore != "" {
		value, err := parseTimestampBound(opts.CreatedBefore, true)
		if err != nil {
			return db.TaskFilter{}, nil, err
		}
		filter.CreatedBefore = &value
	}
	if opts.ModifiedAfter != "" {
		value, err := parseTimestampBound(opts.ModifiedAfter, false)
		if err != nil {
			return db.TaskFilter{}, nil, err
		}
		filter.ModifiedAfter = &value
	}
	if opts.ModifiedBefore != "" {
		value, err := parseTimestampBound(opts.ModifiedBefore, true)
		if err != nil {
			return db.TaskFilter{}, nil, err
		}
		filter.ModifiedBefore = &value
	}
	if opts.DueBefore != "" {
		value, err := parseThingsDate(opts.DueBefore)
		if err != nil {
			return db.TaskFilter{}, nil, err
		}
		filter.DueBefore = &value
	}
	if opts.StartBefore != "" {
		value, err := parseThingsDate(opts.StartBefore)
		if err != nil {
			return db.TaskFilter{}, nil, err
		}
		filter.StartBefore = &value
	}
	if opts.HasURLSet {
		filter.HasURL = &opts.HasURL
	}

	return filter, sortFields, nil
}

func parseTimestampBound(input string, isBefore bool) (float64, error) {
	parsed, dateOnly, err := parseDateOrTime(input)
	if err != nil {
		return 0, err
	}
	if dateOnly && isBefore {
		parsed = parsed.AddDate(0, 0, 1)
	}
	return float64(parsed.Unix()), nil
}

func parseThingsDate(input string) (int, error) {
	parsed, _, err := parseDateOrTime(input)
	if err != nil {
		return 0, err
	}
	return thingsDateValue(parsed.In(time.Local)), nil
}

func parseDateOrTime(input string) (time.Time, bool, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return time.Time{}, false, fmt.Errorf("Error: date required")
	}
	if t, err := time.Parse(time.RFC3339Nano, input); err == nil {
		return t, false, nil
	}
	if t, err := time.Parse(time.RFC3339, input); err == nil {
		return t, false, nil
	}
	if t, err := time.ParseInLocation("2006-01-02 15:04:05", input, time.Local); err == nil {
		return t, false, nil
	}
	if t, err := time.ParseInLocation("2006-01-02 15:04", input, time.Local); err == nil {
		return t, false, nil
	}
	if t, err := time.ParseInLocation("2006-01-02", input, time.Local); err == nil {
		return t, true, nil
	}
	return time.Time{}, false, fmt.Errorf("Error: invalid date %q (use YYYY-MM-DD or RFC3339)", input)
}

func thingsDateValue(t time.Time) int {
	date := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return date.Year()<<16 | int(date.Month())<<12 | date.Day()<<7
}

func parseSortSpec(spec string) ([]TaskSortField, string, error) {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		return nil, "", nil
	}
	parts := strings.Split(spec, ",")
	fields := make([]TaskSortField, 0, len(parts))
	orderParts := make([]string, 0, len(parts))
	for _, raw := range parts {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		desc := false
		if strings.HasPrefix(raw, "-") {
			desc = true
			raw = strings.TrimSpace(strings.TrimPrefix(raw, "-"))
		}
		name := strings.ToLower(raw)
		if alias, ok := taskSortAliases[name]; ok {
			name = alias
		}
		orderExpr, ok := taskSortOrder[name]
		if !ok {
			return nil, "", fmt.Errorf("Error: invalid sort field %q", raw)
		}
		dir := " ASC"
		if desc {
			dir = " DESC"
		}
		orderParts = append(orderParts, orderExpr+dir)
		fields = append(fields, TaskSortField{Field: name, Desc: desc})
	}
	return fields, strings.Join(orderParts, ", "), nil
}

var taskSortAliases = map[string]string{
	"due":         "deadline",
	"proj":        "project",
	"today-index": "today_idx",
	"today_index": "today_idx",
}

var taskSortOrder = map[string]string{
	"created":   "t.creationDate",
	"modified":  "t.userModificationDate",
	"deadline":  "t.deadline",
	"start":     "t.startDate",
	"title":     "t.title COLLATE NOCASE",
	"project":   "p.title COLLATE NOCASE",
	"area":      "a.title COLLATE NOCASE",
	"heading":   "h.title COLLATE NOCASE",
	"status":    "t.status",
	"uuid":      "t.uuid",
	"index":     "t.\"index\"",
	"today_idx": "t.todayIndex",
}
