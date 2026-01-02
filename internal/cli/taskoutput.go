package cli

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/ossianhempel/things3-cli/internal/db"
)

type TaskOutputOptions struct {
	Format   string
	Select   []string
	NoHeader bool
}

func resolveTaskOutputOptions(format string, asJSON bool, selectRaw string, noHeader bool) (TaskOutputOptions, error) {
	format = strings.TrimSpace(strings.ToLower(format))
	if format == "" {
		if asJSON {
			format = "json"
		} else {
			format = "table"
		}
	} else if asJSON && format != "json" {
		return TaskOutputOptions{}, fmt.Errorf("Error: --json cannot be used with --format %s", format)
	}
	switch format {
	case "table", "json", "jsonl", "csv":
	default:
		return TaskOutputOptions{}, fmt.Errorf("Error: invalid format %q", format)
	}
	selectFields, err := parseTaskSelect(selectRaw)
	if err != nil {
		return TaskOutputOptions{}, err
	}
	return TaskOutputOptions{
		Format:   format,
		Select:   selectFields,
		NoHeader: noHeader,
	}, nil
}

func parseTaskSelect(input string) ([]string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, nil
	}
	parts := strings.Split(input, ",")
	fields := make([]string, 0, len(parts))
	seen := map[string]bool{}
	for _, raw := range parts {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		name := normalizeTaskField(raw)
		if name == "" {
			return nil, fmt.Errorf("Error: invalid select field %q (allowed: %s)", raw, strings.Join(sortedTaskFields(), ", "))
		}
		if !seen[name] {
			seen[name] = true
			fields = append(fields, name)
		}
	}
	return fields, nil
}

func normalizeTaskField(raw string) string {
	name := strings.ToLower(strings.TrimSpace(raw))
	if name == "" {
		return ""
	}
	if alias, ok := taskFieldAliases[name]; ok {
		name = alias
	}
	if _, ok := taskFieldHeaders[name]; !ok {
		return ""
	}
	return name
}

var taskFieldAliases = map[string]string{
	"project_title": "project",
	"area_title":    "area",
	"heading_title": "heading",
	"status_label":  "status_label",
	"todayindex":    "today_index",
	"today-index":   "today_index",
	"today_index":   "today_index",
}

var taskFieldHeaders = map[string]string{
	"uuid":         "UUID",
	"title":        "TITLE",
	"project":      "PROJECT",
	"project_id":   "PROJECT_ID",
	"area":         "AREA",
	"area_id":      "AREA_ID",
	"heading":      "HEADING",
	"heading_id":   "HEADING_ID",
	"status":       "STATUS",
	"status_label": "STATUS_LABEL",
	"trashed":      "TRASHED",
	"notes":        "NOTES",
	"start":        "START",
	"start_date":   "START_DATE",
	"deadline":     "DEADLINE",
	"stop_date":    "STOP_DATE",
	"created":      "CREATED",
	"modified":     "MODIFIED",
	"index":        "INDEX",
	"today_index":  "TODAY_INDEX",
	"tags":         "TAGS",
	"type":         "TYPE",
}

var defaultTaskTableFields = []string{
	"uuid",
	"title",
	"project",
	"area",
	"heading",
	"status",
	"trashed",
}

func taskFieldValue(task db.Task, field string) any {
	switch field {
	case "uuid":
		return task.UUID
	case "title":
		return task.Title
	case "project":
		return task.ProjectTitle
	case "project_id":
		return task.ProjectID
	case "area":
		return task.AreaTitle
	case "area_id":
		return task.AreaID
	case "heading":
		return task.HeadingTitle
	case "heading_id":
		return task.HeadingID
	case "status":
		return task.Status
	case "status_label":
		return db.StatusLabel(task.Status)
	case "trashed":
		return task.Trashed
	case "notes":
		return task.Notes
	case "start":
		return task.Start
	case "start_date":
		return task.StartDate
	case "deadline":
		return task.Deadline
	case "stop_date":
		return task.StopDate
	case "created":
		return task.Created
	case "modified":
		return task.Modified
	case "index":
		return task.Index
	case "today_index":
		if task.TodayIndex == nil {
			return nil
		}
		return *task.TodayIndex
	case "tags":
		return task.Tags
	case "type":
		return task.Type
	default:
		return ""
	}
}

func taskFieldString(task db.Task, field string) string {
	switch field {
	case "status":
		return db.StatusLabel(task.Status)
	case "status_label":
		return db.StatusLabel(task.Status)
	case "trashed":
		return strconv.FormatBool(task.Trashed)
	case "index":
		if task.Index == 0 {
			return ""
		}
		return strconv.Itoa(task.Index)
	case "today_index":
		if task.TodayIndex == nil {
			return ""
		}
		return strconv.Itoa(*task.TodayIndex)
	case "tags":
		return strings.Join(task.Tags, ",")
	default:
		value := taskFieldValue(task, field)
		switch v := value.(type) {
		case nil:
			return ""
		case string:
			return v
		case int:
			return strconv.Itoa(v)
		case bool:
			return strconv.FormatBool(v)
		case []string:
			return strings.Join(v, ",")
		default:
			return fmt.Sprintf("%v", v)
		}
	}
}

func writeTasks(out io.Writer, tasks []db.Task, opts TaskOutputOptions) error {
	switch opts.Format {
	case "json":
		enc := json.NewEncoder(out)
		if len(opts.Select) == 0 {
			return enc.Encode(tasks)
		}
		records := make([]map[string]any, 0, len(tasks))
		for _, task := range tasks {
			record := make(map[string]any, len(opts.Select))
			for _, field := range opts.Select {
				record[field] = taskFieldValue(task, field)
			}
			records = append(records, record)
		}
		return enc.Encode(records)
	case "jsonl":
		enc := json.NewEncoder(out)
		if len(opts.Select) == 0 {
			for _, task := range tasks {
				if err := enc.Encode(task); err != nil {
					return err
				}
			}
			return nil
		}
		for _, task := range tasks {
			record := make(map[string]any, len(opts.Select))
			for _, field := range opts.Select {
				record[field] = taskFieldValue(task, field)
			}
			if err := enc.Encode(record); err != nil {
				return err
			}
		}
		return nil
	case "csv":
		fields := opts.Select
		if len(fields) == 0 {
			fields = defaultTaskTableFields
		}
		writer := csv.NewWriter(out)
		if !opts.NoHeader {
			headers := make([]string, 0, len(fields))
			for _, field := range fields {
				headers = append(headers, taskFieldHeaders[field])
			}
			if err := writer.Write(headers); err != nil {
				return err
			}
		}
		for _, task := range tasks {
			row := make([]string, 0, len(fields))
			for _, field := range fields {
				row = append(row, taskFieldString(task, field))
			}
			if err := writer.Write(row); err != nil {
				return err
			}
		}
		writer.Flush()
		return writer.Error()
	case "table":
		fields := opts.Select
		if len(fields) == 0 {
			fields = defaultTaskTableFields
		}
		return writeTaskTable(out, tasks, fields, opts.NoHeader)
	default:
		return fmt.Errorf("Error: invalid format %q", opts.Format)
	}
}

func writeTaskTable(out io.Writer, tasks []db.Task, fields []string, noHeader bool) error {
	w := tabwriter.NewWriter(out, 0, 2, 2, ' ', 0)
	if !noHeader {
		headers := make([]string, 0, len(fields))
		for _, field := range fields {
			headers = append(headers, taskFieldHeaders[field])
		}
		fmt.Fprintln(w, strings.Join(headers, "\t"))
	}
	for _, task := range tasks {
		row := make([]string, 0, len(fields))
		for _, field := range fields {
			row = append(row, taskFieldString(task, field))
		}
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}
	return w.Flush()
}

func sortedTaskFields() []string {
	fields := make([]string, 0, len(taskFieldHeaders))
	for field := range taskFieldHeaders {
		fields = append(fields, field)
	}
	sort.Strings(fields)
	return fields
}
