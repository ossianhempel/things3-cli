package cli

import (
	"sort"
	"strings"

	"github.com/ossianhempel/things3-cli/internal/db"
)

func sortTasks(tasks []db.Task, spec []TaskSortField) {
	if len(spec) == 0 || len(tasks) < 2 {
		return
	}
	sort.SliceStable(tasks, func(i, j int) bool {
		left := tasks[i]
		right := tasks[j]
		for _, field := range spec {
			cmp := compareTaskField(left, right, field.Field)
			if cmp == 0 {
				continue
			}
			if field.Desc {
				return cmp > 0
			}
			return cmp < 0
		}
		return left.UUID < right.UUID
	})
}

func compareTaskField(left db.Task, right db.Task, field string) int {
	switch field {
	case "created":
		return compareString(left.Created, right.Created)
	case "modified":
		return compareString(left.Modified, right.Modified)
	case "deadline":
		return compareString(left.Deadline, right.Deadline)
	case "start":
		return compareString(left.StartDate, right.StartDate)
	case "title":
		return compareStringCI(left.Title, right.Title)
	case "project":
		return compareStringCI(left.ProjectTitle, right.ProjectTitle)
	case "area":
		return compareStringCI(left.AreaTitle, right.AreaTitle)
	case "heading":
		return compareStringCI(left.HeadingTitle, right.HeadingTitle)
	case "status":
		return compareInt(left.Status, right.Status)
	case "uuid":
		return compareString(left.UUID, right.UUID)
	case "index":
		return compareInt(left.Index, right.Index)
	case "today_idx":
		return compareInt(ptrToInt(left.TodayIndex), ptrToInt(right.TodayIndex))
	default:
		return 0
	}
}

func compareString(left string, right string) int {
	if left == "" && right == "" {
		return 0
	}
	if left == "" {
		return 1
	}
	if right == "" {
		return -1
	}
	if left < right {
		return -1
	}
	if left > right {
		return 1
	}
	return 0
}

func compareStringCI(left string, right string) int {
	return compareString(strings.ToLower(left), strings.ToLower(right))
}

func compareInt(left int, right int) int {
	if left < right {
		return -1
	}
	if left > right {
		return 1
	}
	return 0
}

func ptrToInt(value *int) int {
	if value == nil {
		return 0
	}
	return *value
}
