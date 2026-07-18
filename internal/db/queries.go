package db

import (
	"cmp"
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const tagSeparator = "\x1f"

// Projects returns projects in the database.
func (s *Store) Projects(filter ProjectFilter) ([]Project, error) {
	if s == nil || s.conn == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	var b strings.Builder
	b.WriteString("SELECT t.uuid, t.title, t.status, t.trashed, t.area, a.title ")
	b.WriteString("FROM TMTask t ")
	b.WriteString("LEFT JOIN TMArea a ON t.area = a.uuid ")
	b.WriteString("WHERE t.type = ?")
	args := []any{TaskTypeProject}

	if filter.AreaID != "" {
		b.WriteString(" AND t.area = ?")
		args = append(args, filter.AreaID)
	}
	if !filter.IncludeTrashed {
		b.WriteString(" AND t.trashed = 0")
	}
	if filter.Status != nil {
		b.WriteString(" AND t.status = ?")
		args = append(args, *filter.Status)
	}
	b.WriteString(" AND t.rt1_recurrenceRule IS NULL")
	b.WriteString(" ORDER BY t.\"index\"")

	rows, err := s.conn.Query(b.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]Project, 0, 64)
	for rows.Next() {
		var p Project
		var areaID sql.NullString
		var areaTitle sql.NullString
		if err := rows.Scan(&p.UUID, &p.Title, &p.Status, &p.Trashed, &areaID, &areaTitle); err != nil {
			return nil, err
		}
		if areaID.Valid {
			p.AreaID = areaID.String
		}
		if areaTitle.Valid {
			p.AreaTitle = areaTitle.String
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

// ProjectsWithoutArea returns projects that are not assigned to an area.
func (s *Store) ProjectsWithoutArea(filter ProjectFilter) ([]Project, error) {
	if s == nil || s.conn == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	var b strings.Builder
	b.WriteString("SELECT t.uuid, t.title, t.status, t.trashed, t.area, a.title ")
	b.WriteString("FROM TMTask t ")
	b.WriteString("LEFT JOIN TMArea a ON t.area = a.uuid ")
	b.WriteString("WHERE t.type = ? AND t.area IS NULL")
	args := []any{TaskTypeProject}

	if !filter.IncludeTrashed {
		b.WriteString(" AND t.trashed = 0")
	}
	if filter.Status != nil {
		b.WriteString(" AND t.status = ?")
		args = append(args, *filter.Status)
	}
	b.WriteString(" AND t.rt1_recurrenceRule IS NULL")
	b.WriteString(" ORDER BY t.\"index\"")

	rows, err := s.conn.Query(b.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]Project, 0, 32)
	for rows.Next() {
		var p Project
		var areaID sql.NullString
		var areaTitle sql.NullString
		if err := rows.Scan(&p.UUID, &p.Title, &p.Status, &p.Trashed, &areaID, &areaTitle); err != nil {
			return nil, err
		}
		if areaID.Valid {
			p.AreaID = areaID.String
		}
		if areaTitle.Valid {
			p.AreaTitle = areaTitle.String
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

// Areas returns areas in the database.
func (s *Store) Areas() ([]Area, error) {
	if s == nil || s.conn == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	rows, err := s.conn.Query("SELECT uuid, title, visible FROM TMArea ORDER BY `index`")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	areas := make([]Area, 0, 32)
	for rows.Next() {
		var a Area
		var visible sql.NullInt64
		if err := rows.Scan(&a.UUID, &a.Title, &visible); err != nil {
			return nil, err
		}
		if visible.Valid {
			a.Visible = visible.Int64 != 0
		}
		areas = append(areas, a)
	}
	return areas, rows.Err()
}

// Tags returns tags in the database.
func (s *Store) Tags() ([]Tag, error) {
	if s == nil || s.conn == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	rows, err := s.conn.Query(
		`SELECT t.uuid, t.title, t.shortcut, t.parent, COUNT(tt.tasks) AS usage
		 FROM TMTag t
		 LEFT JOIN TMTaskTag tt ON tt.tags = t.uuid
		 GROUP BY t.uuid
		 ORDER BY usage DESC, t.title COLLATE NOCASE`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make([]Tag, 0, 64)
	for rows.Next() {
		var t Tag
		var shortcut sql.NullString
		var parent sql.NullString
		if err := rows.Scan(&t.UUID, &t.Title, &shortcut, &parent, &t.Usage); err != nil {
			return nil, err
		}
		if shortcut.Valid {
			t.Shortcut = shortcut.String
		}
		if parent.Valid {
			t.ParentID = parent.String
		}
		tags = append(tags, t)
	}
	return tags, rows.Err()
}

// Tasks returns tasks in the database.
func (s *Store) Tasks(filter TaskFilter) ([]Task, error) {
	return s.queryTasks("", nil, filter, "")
}

// TemplatesTasks returns hidden repeating template tasks in the database.
func (s *Store) TemplatesTasks(filter TaskFilter) ([]Task, error) {
	filter.IncludeRepeating = true
	where := "t.rt1_repeatingTemplate IS NULL AND (t.rt1_recurrenceRule IS NOT NULL OR t.repeater IS NOT NULL)"
	tasks, err := s.queryTasks(where, nil, filter, "")
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(tasks))
	for i := range tasks {
		ids[i] = tasks[i].UUID
	}
	states, err := s.RepeatStatesByIDs(ids)
	if err != nil {
		return nil, err
	}
	for i := range tasks {
		tasks[i].Repeat = states[tasks[i].UUID]
	}
	return tasks, nil
}

// TaskByID returns a single task by UUID.
func (s *Store) TaskByID(id string) (*Task, error) {
	if strings.TrimSpace(id) == "" {
		return nil, sql.ErrNoRows
	}
	filter := TaskFilter{
		IncludeTrashed:        true,
		ExcludeTrashedContext: false,
		IncludeChecklist:      true,
	}
	tasks, err := s.queryTasks("t.uuid = ?", []any{id}, filter, "")
	if err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return nil, sql.ErrNoRows
	}
	return &tasks[0], nil
}

// TodayTasks returns tasks that belong in Today according to Things rules.
func (s *Store) TodayTasks(filter TaskFilter) ([]Task, error) {
	todayExpr := thingsDateTodayExpr()
	const todayOrder = "t.startBucket ASC, t.todayIndexReferenceDate DESC, t.todayIndex ASC, t.uuid ASC"
	regular, err := s.queryTasks("t.start = 1 AND t.startDate IS NOT NULL", nil, filter, todayOrder)
	if err != nil {
		return nil, err
	}
	unconfirmedScheduled, err := s.queryTasks("t.start = 2 AND t.startDate IS NOT NULL AND t.startDate <= "+todayExpr, nil, filter, todayOrder)
	if err != nil {
		return nil, err
	}
	unconfirmedOverdue, err := s.queryTasks("t.startDate IS NULL AND t.deadline IS NOT NULL AND t.deadline <= "+todayExpr+" AND t.deadlineSuppressionDate IS NULL", nil, filter, todayOrder)
	if err != nil {
		return nil, err
	}
	result := append(regular, unconfirmedScheduled...)
	result = append(result, unconfirmedOverdue...)
	sort.Slice(result, func(i, j int) bool {
		return todayTaskLess(result[i], result[j])
	})
	return result, nil
}

func todayTaskLess(left, right Task) bool {
	if cmp := compareNullableInt(left.StartBucket, right.StartBucket, false); cmp != 0 {
		return cmp < 0
	}
	if cmp := compareNullableInt(left.TodayIndexReferenceDate, right.TodayIndexReferenceDate, true); cmp != 0 {
		return cmp < 0
	}
	if cmp := compareNullableInt(left.TodayIndex, right.TodayIndex, false); cmp != 0 {
		return cmp < 0
	}
	return left.UUID < right.UUID
}

// compareNullableInt matches SQLite's native null placement: NULL sorts first
// in ascending order and last in descending order.
func compareNullableInt(left, right *int, descending bool) int {
	if left == nil && right == nil {
		return 0
	}
	if left == nil {
		if descending {
			return 1
		}
		return -1
	}
	if right == nil {
		if descending {
			return -1
		}
		return 1
	}
	comparison := cmp.Compare(*left, *right)
	if descending {
		return -comparison
	}
	return comparison
}

// InboxTasks returns inbox tasks.
func (s *Store) InboxTasks(filter TaskFilter) ([]Task, error) {
	return s.queryTasks("t.start = 0", nil, filter, "")
}

// AnytimeTasks returns Anytime tasks.
func (s *Store) AnytimeTasks(filter TaskFilter) ([]Task, error) {
	return s.queryTasks("t.start = 1", nil, filter, "")
}

// SomedayTasks returns Someday tasks.
func (s *Store) SomedayTasks(filter TaskFilter) ([]Task, error) {
	return s.queryTasks("t.start = 2 AND t.startDate IS NULL", nil, filter, "")
}

// UpcomingTasks returns upcoming tasks (scheduled or with upcoming deadlines).
func (s *Store) UpcomingTasks(filter TaskFilter) ([]Task, error) {
	todayExpr := thingsDateTodayExpr()
	where := "t.start = 2 AND t.startDate IS NOT NULL AND t.startDate > " + todayExpr
	return s.queryTasks(where, nil, filter, "")
}

// DeadlinesTasks returns tasks with deadlines.
func (s *Store) DeadlinesTasks(filter TaskFilter) ([]Task, error) {
	tasks, err := s.queryTasks("t.deadline IS NOT NULL", nil, filter, "")
	if err != nil {
		return nil, err
	}
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Deadline < tasks[j].Deadline
	})
	return tasks, nil
}

// LogbookTasks returns completed or canceled tasks.
func (s *Store) LogbookTasks(filter TaskFilter) ([]Task, error) {
	where := "t.status IN (?, ?)"
	args := []any{StatusCompleted, StatusCanceled}
	order := "t.stopDate DESC"
	return s.queryTasks(where, args, filter, order)
}

// CompletedTasks returns completed tasks.
func (s *Store) CompletedTasks(filter TaskFilter) ([]Task, error) {
	status := StatusCompleted
	filter.Status = &status
	return s.queryTasks("", nil, filter, "")
}

// CanceledTasks returns canceled tasks.
func (s *Store) CanceledTasks(filter TaskFilter) ([]Task, error) {
	status := StatusCanceled
	filter.Status = &status
	return s.queryTasks("", nil, filter, "")
}

// TrashTasks returns trashed tasks.
func (s *Store) TrashTasks(filter TaskFilter) ([]Task, error) {
	filter.IncludeTrashed = true
	filter.ExcludeTrashedContext = false
	where := "t.trashed = 1"
	return s.queryTasks(where, nil, filter, "")
}

// TasksCreatedBetween returns tasks created within the given range.
func (s *Store) TasksCreatedBetween(start, end time.Time, filter TaskFilter) ([]Task, error) {
	where := "t.creationDate >= ? AND t.creationDate <= ?"
	args := []any{float64(start.Unix()), float64(end.Unix())}
	order := "t.creationDate DESC"
	return s.queryTasks(where, args, filter, order)
}

// TasksCompletedBetween returns completed or canceled tasks within the given range.
func (s *Store) TasksCompletedBetween(start, end time.Time, filter TaskFilter) ([]Task, error) {
	where := "t.stopDate >= ? AND t.stopDate < ? AND t.status IN (?, ?)"
	args := []any{float64(start.Unix()), float64(end.Unix()), StatusCompleted, StatusCanceled}
	order := "t.stopDate DESC"
	return s.queryTasks(where, args, filter, order)
}

// ItemByID returns a single item (task/project/heading, area, or tag) by UUID.
func (s *Store) ItemByID(id string, includeChecklist bool) (*Item, error) {
	if s == nil || s.conn == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	row := s.conn.QueryRow(
		`SELECT t.uuid, t.type, t.title, t.status, t.trashed, p.title, a.title, h.title
         FROM TMTask t
         LEFT JOIN TMTask p ON t.project = p.uuid
         LEFT JOIN TMArea a ON t.area = a.uuid
         LEFT JOIN TMTask h ON t.heading = h.uuid
         WHERE t.uuid = ?`,
		id,
	)
	var item Item
	var taskType int
	var projectTitle sql.NullString
	var areaTitle sql.NullString
	var headingTitle sql.NullString
	var status int
	var trashed int
	if err := row.Scan(&item.UUID, &taskType, &item.Title, &status, &trashed, &projectTitle, &areaTitle, &headingTitle); err == nil {
		item.Type = taskTypeLabel(taskType)
		item.Status = &status
		trashedBool := trashed != 0
		item.Trashed = &trashedBool
		if projectTitle.Valid {
			item.ProjectTitle = projectTitle.String
		}
		if areaTitle.Valid {
			item.AreaTitle = areaTitle.String
		}
		if headingTitle.Valid {
			item.HeadingTitle = headingTitle.String
		}
		if includeChecklist && item.Type == "to-do" {
			checklist, err := loadChecklistItems(s.conn, []string{item.UUID})
			if err != nil {
				return nil, err
			}
			item.Checklist = checklist[item.UUID]
		}
		return &item, nil
	} else if err != sql.ErrNoRows {
		return nil, err
	}

	row = s.conn.QueryRow("SELECT uuid, title, visible FROM TMArea WHERE uuid = ?", id)
	var visible sql.NullInt64
	if err := row.Scan(&item.UUID, &item.Title, &visible); err == nil {
		item.Type = "area"
		if visible.Valid {
			visibleBool := visible.Int64 != 0
			item.Visible = &visibleBool
		}
		return &item, nil
	} else if err != sql.ErrNoRows {
		return nil, err
	}

	row = s.conn.QueryRow("SELECT uuid, title, shortcut, parent FROM TMTag WHERE uuid = ?", id)
	var shortcut sql.NullString
	var parent sql.NullString
	if err := row.Scan(&item.UUID, &item.Title, &shortcut, &parent); err == nil {
		item.Type = "tag"
		if shortcut.Valid {
			item.Shortcut = shortcut.String
		}
		if parent.Valid {
			item.ParentID = parent.String
		}
		return &item, nil
	} else if err != sql.ErrNoRows {
		return nil, err
	}

	return nil, sql.ErrNoRows
}

// ItemsByTitle returns matching items by exact title (case-insensitive).
func (s *Store) ItemsByTitle(title string, includeChecklist bool) ([]Item, error) {
	if s == nil || s.conn == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	items := make([]Item, 0, 4)

	rows, err := s.conn.Query(
		`SELECT t.uuid, t.type, t.title, t.status, t.trashed, p.title, a.title, h.title
         FROM TMTask t
         LEFT JOIN TMTask p ON t.project = p.uuid
         LEFT JOIN TMArea a ON t.area = a.uuid
         LEFT JOIN TMTask h ON t.heading = h.uuid
         WHERE lower(t.title) = lower(?)`,
		title,
	)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var item Item
		var taskType int
		var projectTitle sql.NullString
		var areaTitle sql.NullString
		var headingTitle sql.NullString
		var status int
		var trashed int
		if err := rows.Scan(&item.UUID, &taskType, &item.Title, &status, &trashed, &projectTitle, &areaTitle, &headingTitle); err != nil {
			rows.Close()
			return nil, err
		}
		item.Type = taskTypeLabel(taskType)
		item.Status = &status
		trashedBool := trashed != 0
		item.Trashed = &trashedBool
		if projectTitle.Valid {
			item.ProjectTitle = projectTitle.String
		}
		if areaTitle.Valid {
			item.AreaTitle = areaTitle.String
		}
		if headingTitle.Valid {
			item.HeadingTitle = headingTitle.String
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}
	rows.Close()

	if includeChecklist {
		taskIDs := make([]string, 0, len(items))
		for _, item := range items {
			if item.Type == "to-do" {
				taskIDs = append(taskIDs, item.UUID)
			}
		}
		if len(taskIDs) > 0 {
			checklist, err := loadChecklistItems(s.conn, taskIDs)
			if err != nil {
				return nil, err
			}
			for i := range items {
				if items[i].Type == "to-do" {
					items[i].Checklist = checklist[items[i].UUID]
				}
			}
		}
	}

	rows, err = s.conn.Query("SELECT uuid, title, visible FROM TMArea WHERE lower(title) = lower(?)", title)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var item Item
		var visible sql.NullInt64
		if err := rows.Scan(&item.UUID, &item.Title, &visible); err != nil {
			rows.Close()
			return nil, err
		}
		item.Type = "area"
		if visible.Valid {
			visibleBool := visible.Int64 != 0
			item.Visible = &visibleBool
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}
	rows.Close()

	rows, err = s.conn.Query("SELECT uuid, title, shortcut, parent FROM TMTag WHERE lower(title) = lower(?)", title)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var item Item
		var shortcut sql.NullString
		var parent sql.NullString
		if err := rows.Scan(&item.UUID, &item.Title, &shortcut, &parent); err != nil {
			rows.Close()
			return nil, err
		}
		item.Type = "tag"
		if shortcut.Valid {
			item.Shortcut = shortcut.String
		}
		if parent.Valid {
			item.ParentID = parent.String
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}
	rows.Close()

	return items, nil
}

func taskTypeLabel(taskType int) string {
	switch taskType {
	case TaskTypeTodo:
		return "to-do"
	case TaskTypeProject:
		return "project"
	case TaskTypeHeading:
		return "heading"
	default:
		return "task"
	}
}

// ResolveAreaID resolves an area by UUID or title.
func (s *Store) ResolveAreaID(input string) (string, error) {
	return resolveAreaID(s.conn, input)
}

// ResolveProjectID resolves a project by UUID or title.
func (s *Store) ResolveProjectID(input string) (string, error) {
	return resolveProjectID(s.conn, input)
}

// ResolveTagID resolves a tag by UUID or title.
func (s *Store) ResolveTagID(input string) (string, error) {
	return resolveTagID(s.conn, input)
}

func resolveAreaID(conn *sql.DB, input string) (string, error) {
	if input == "" {
		return "", nil
	}
	var id string
	if err := conn.QueryRow("SELECT uuid FROM TMArea WHERE uuid = ?", input).Scan(&id); err == nil {
		return id, nil
	} else if err != sql.ErrNoRows {
		return "", err
	}
	if err := conn.QueryRow("SELECT uuid FROM TMArea WHERE lower(title) = lower(?)", input).Scan(&id); err == nil {
		return id, nil
	} else if err != sql.ErrNoRows {
		return "", err
	}
	return resolveByNameFallback(conn, input, "area", "SELECT uuid, title FROM TMArea")
}

func resolveProjectID(conn *sql.DB, input string) (string, error) {
	if input == "" {
		return "", nil
	}
	var id string
	if err := conn.QueryRow("SELECT uuid FROM TMTask WHERE type = ? AND uuid = ?", TaskTypeProject, input).Scan(&id); err == nil {
		return id, nil
	} else if err != sql.ErrNoRows {
		return "", err
	}
	if err := conn.QueryRow("SELECT uuid FROM TMTask WHERE type = ? AND lower(title) = lower(?)", TaskTypeProject, input).Scan(&id); err == nil {
		return id, nil
	} else if err != sql.ErrNoRows {
		return "", err
	}
	// Exclude trashed projects from fuzzy matching: a loose substring match
	// should never silently resolve to a trashed project on a write path. The
	// exact UUID/title tiers above intentionally still resolve trashed projects,
	// since those express explicit intent.
	return resolveByNameFallback(conn, input, "project", "SELECT uuid, title FROM TMTask WHERE type = ? AND trashed = 0", TaskTypeProject)
}

func resolveTagID(conn *sql.DB, input string) (string, error) {
	if input == "" {
		return "", nil
	}
	var id string
	if err := conn.QueryRow("SELECT uuid FROM TMTag WHERE uuid = ?", input).Scan(&id); err == nil {
		return id, nil
	} else if err != sql.ErrNoRows {
		return "", err
	}
	if err := conn.QueryRow("SELECT uuid FROM TMTag WHERE lower(title) = lower(?)", input).Scan(&id); err == nil {
		return id, nil
	} else if err != sql.ErrNoRows {
		return "", err
	}
	return resolveByNameFallback(conn, input, "tag", "SELECT uuid, title FROM TMTag")
}

// normalizeForMatch lowercases s and keeps only letters, digits, and single
// interior spaces, dropping emoji, symbols, and punctuation. This powers
// emoji-insensitive name resolution so a query like "Rebtech" can match an
// area titled "💼 Rebtech".
func normalizeForMatch(s string) string {
	var b strings.Builder
	lastSpace := false
	for _, r := range strings.ToLower(s) {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
			lastSpace = false
		case unicode.IsSpace(r):
			if b.Len() > 0 && !lastSpace {
				b.WriteRune(' ')
				lastSpace = true
			}
		}
		// Everything else (emoji, punctuation, symbols) is dropped.
	}
	return strings.TrimRight(b.String(), " ")
}

// resolveByNameFallback is the loose-matching tier used after exact UUID and
// exact title lookups fail. It compares an emoji/symbol-stripped form of the
// input against every candidate title, preferring a normalized-exact match and
// falling back to a substring match. It resolves only when exactly one
// candidate matches; multiple matches return an error listing them so the
// caller can disambiguate. listQuery must select (uuid, title) for the entity.
func resolveByNameFallback(conn *sql.DB, input, label, listQuery string, listArgs ...any) (string, error) {
	normInput := normalizeForMatch(input)
	if normInput == "" {
		return "", fmt.Errorf("%s not found: %s", label, input)
	}

	rows, err := conn.Query(listQuery, listArgs...)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	type candidate struct {
		uuid  string
		title string
	}
	var exact, substr []candidate
	for rows.Next() {
		var uuid string
		var title sql.NullString
		if err := rows.Scan(&uuid, &title); err != nil {
			return "", err
		}
		norm := normalizeForMatch(title.String)
		if norm == "" {
			continue
		}
		switch {
		case norm == normInput:
			exact = append(exact, candidate{uuid, title.String})
		case strings.Contains(norm, normInput):
			substr = append(substr, candidate{uuid, title.String})
		}
	}
	if err := rows.Err(); err != nil {
		return "", err
	}

	matches := exact
	if len(matches) == 0 {
		matches = substr
	}
	switch len(matches) {
	case 1:
		return matches[0].uuid, nil
	case 0:
		return "", fmt.Errorf("%s not found: %s", label, input)
	default:
		titles := make([]string, len(matches))
		for i, m := range matches {
			titles[i] = strconv.Quote(m.title)
		}
		return "", fmt.Errorf("%s %q is ambiguous, matches: %s", label, input, strings.Join(titles, ", "))
	}
}

func thingsDateTodayExpr() string {
	return "((strftime('%Y', date('now', 'localtime')) << 16) | (strftime('%m', date('now', 'localtime')) << 12) | (strftime('%d', date('now', 'localtime')) << 7))"
}

// queryTasks runs the common task query with the given WHERE clause and filter.
// SECURITY: where must only contain hardcoded SQL fragments with ? placeholders
// for parameter binding — never pass unsanitized user input as the where string.
func (s *Store) queryTasks(where string, args []any, filter TaskFilter, order string) ([]Task, error) {
	if s == nil || s.conn == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	var b strings.Builder
	b.WriteString("SELECT t.uuid, t.type, t.title, t.status, t.trashed, t.notes, t.start, t.startDate, t.startBucket, t.deadline, t.stopDate, t.creationDate, t.userModificationDate, t.\"index\", t.todayIndex, t.todayIndexReferenceDate, (t.rt1_recurrenceRule IS NOT NULL OR t.repeater IS NOT NULL) AS repeating, ")
	b.WriteString("t.project, p.title, t.area, a.title, t.heading, h.title, ")
	b.WriteString("(SELECT group_concat(title, '" + tagSeparator + "') FROM (")
	b.WriteString("SELECT tag.title AS title FROM TMTag tag ")
	b.WriteString("JOIN TMTaskTag tt ON tt.tags = tag.uuid ")
	b.WriteString("WHERE tt.tasks = t.uuid ORDER BY tag.title COLLATE NOCASE")
	b.WriteString(")) ")
	b.WriteString("FROM TMTask t ")
	b.WriteString("LEFT JOIN TMTask p ON t.project = p.uuid ")
	b.WriteString("LEFT JOIN TMArea a ON t.area = a.uuid ")
	b.WriteString("LEFT JOIN TMTask h ON t.heading = h.uuid ")
	b.WriteString("LEFT JOIN TMTask hp ON h.project = hp.uuid ")
	b.WriteString("WHERE 1=1")

	params := []any{}
	if len(filter.Types) > 0 {
		placeholders := strings.Repeat("?,", len(filter.Types))
		placeholders = strings.TrimRight(placeholders, ",")
		b.WriteString(" AND t.type IN (" + placeholders + ")")
		for _, t := range filter.Types {
			params = append(params, t)
		}
	}
	if where != "" {
		b.WriteString(" AND (" + where + ")")
		params = append(params, args...)
	}
	if !filter.IncludeTrashed {
		b.WriteString(" AND t.trashed = 0")
	}
	if filter.ExcludeTrashedContext {
		b.WriteString(" AND NOT IFNULL(p.trashed, 0)")
		b.WriteString(" AND NOT IFNULL(hp.trashed, 0)")
	}
	if filter.Status != nil {
		b.WriteString(" AND t.status = ?")
		params = append(params, *filter.Status)
	}
	if filter.ProjectID != "" {
		b.WriteString(" AND (t.uuid = ? OR t.project = ? OR hp.uuid = ?)")
		params = append(params, filter.ProjectID, filter.ProjectID, filter.ProjectID)
	}
	if filter.AreaID != "" {
		b.WriteString(" AND t.area = ?")
		params = append(params, filter.AreaID)
	}
	if filter.TagID != "" {
		b.WriteString(" AND EXISTS (SELECT 1 FROM TMTaskTag tt WHERE tt.tasks = t.uuid AND tt.tags = ?)")
		params = append(params, filter.TagID)
	}
	if filter.Search != "" {
		b.WriteString(` AND (lower(t.title) LIKE lower(?) ESCAPE '\' OR lower(t.notes) LIKE lower(?) ESCAPE '\' OR lower(a.title) LIKE lower(?) ESCAPE '\')`)
		like := "%" + escapeLike(filter.Search) + "%"
		params = append(params, like, like, like)
	}
	if filter.CreatedAfter != nil {
		b.WriteString(" AND t.creationDate >= ?")
		params = append(params, *filter.CreatedAfter)
	}
	if filter.CreatedBefore != nil {
		b.WriteString(" AND t.creationDate < ?")
		params = append(params, *filter.CreatedBefore)
	}
	if filter.ModifiedAfter != nil {
		b.WriteString(" AND t.userModificationDate >= ?")
		params = append(params, *filter.ModifiedAfter)
	}
	if filter.ModifiedBefore != nil {
		b.WriteString(" AND t.userModificationDate < ?")
		params = append(params, *filter.ModifiedBefore)
	}
	if filter.DueBefore != nil {
		b.WriteString(" AND t.deadline IS NOT NULL AND t.deadline <= ?")
		params = append(params, *filter.DueBefore)
	}
	if filter.StartBefore != nil {
		b.WriteString(" AND t.startDate IS NOT NULL AND t.startDate <= ?")
		params = append(params, *filter.StartBefore)
	}
	if filter.HasURL != nil {
		if *filter.HasURL {
			b.WriteString(" AND (IFNULL(t.notes, '') LIKE '%http://%' OR IFNULL(t.notes, '') LIKE '%https://%')")
		} else {
			b.WriteString(" AND (IFNULL(t.notes, '') NOT LIKE '%http://%' AND IFNULL(t.notes, '') NOT LIKE '%https://%')")
		}
	}
	if filter.RepeatingOnly {
		b.WriteString(" AND (t.rt1_recurrenceRule IS NOT NULL OR t.repeater IS NOT NULL)")
	} else if !filter.IncludeRepeating {
		b.WriteString(" AND t.rt1_recurrenceRule IS NULL AND t.repeater IS NULL")
	}

	orderClause := order
	if filter.Order != "" {
		orderClause = filter.Order
	}
	if orderClause == "" {
		orderClause = "t.\"index\""
	}
	b.WriteString(" ORDER BY " + orderClause)
	if filter.Limit > 0 {
		b.WriteString(" LIMIT ?")
		params = append(params, filter.Limit)
		if filter.Offset > 0 {
			b.WriteString(" OFFSET ?")
			params = append(params, filter.Offset)
		}
	} else if filter.Offset > 0 {
		b.WriteString(" LIMIT -1 OFFSET ?")
		params = append(params, filter.Offset)
	}

	rows, err := s.conn.Query(b.String(), params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks, err := scanTaskRows(rows)
	if err != nil {
		return nil, err
	}
	if filter.IncludeChecklist && len(tasks) > 0 {
		ids := make([]string, len(tasks))
		for i, task := range tasks {
			ids[i] = task.UUID
		}
		checklist, err := loadChecklistItems(s.conn, ids)
		if err != nil {
			return nil, err
		}
		for i := range tasks {
			if items, ok := checklist[tasks[i].UUID]; ok {
				tasks[i].Checklist = items
			}
		}
	}
	return tasks, nil
}

func scanTaskRows(rows *sql.Rows) ([]Task, error) {
	tasks := make([]Task, 0, 128)
	for rows.Next() {
		var t Task
		var taskType int
		var notes sql.NullString
		var start sql.NullInt64
		var startDate sql.NullInt64
		var startBucket sql.NullInt64
		var deadline sql.NullInt64
		var stopDate sql.NullFloat64
		var created sql.NullFloat64
		var modified sql.NullFloat64
		var index sql.NullInt64
		var todayIndex sql.NullInt64
		var todayIndexReferenceDate sql.NullInt64
		var repeating sql.NullInt64
		var projectID sql.NullString
		var projectTitle sql.NullString
		var areaID sql.NullString
		var areaTitle sql.NullString
		var headingID sql.NullString
		var headingTitle sql.NullString
		var tagTitles sql.NullString
		if err := rows.Scan(&t.UUID, &taskType, &t.Title, &t.Status, &t.Trashed, &notes, &start, &startDate, &startBucket, &deadline, &stopDate, &created, &modified, &index, &todayIndex, &todayIndexReferenceDate, &repeating, &projectID, &projectTitle, &areaID, &areaTitle, &headingID, &headingTitle, &tagTitles); err != nil {
			return nil, err
		}
		t.Type = taskTypeLabel(taskType)
		if index.Valid {
			t.Index = int(index.Int64)
		}
		if todayIndex.Valid {
			val := int(todayIndex.Int64)
			t.TodayIndex = &val
		}
		if todayIndexReferenceDate.Valid {
			val := int(todayIndexReferenceDate.Int64)
			t.TodayIndexReferenceDate = &val
		}
		if repeating.Valid {
			t.Repeating = repeating.Int64 != 0
		}
		if notes.Valid {
			t.Notes = notes.String
		}
		if start.Valid {
			t.Start = startLabel(int(start.Int64))
		}
		if startDate.Valid {
			t.StartDate = formatThingsDate(startDate.Int64)
		}
		if startBucket.Valid {
			val := int(startBucket.Int64)
			t.StartBucket = &val
		}
		if deadline.Valid {
			t.Deadline = formatThingsDate(deadline.Int64)
		}
		if stopDate.Valid {
			t.StopDate = formatTimestamp(stopDate.Float64)
		}
		if created.Valid {
			t.Created = formatTimestamp(created.Float64)
		}
		if modified.Valid {
			t.Modified = formatTimestamp(modified.Float64)
		}
		if projectID.Valid {
			t.ProjectID = projectID.String
		}
		if projectTitle.Valid {
			t.ProjectTitle = projectTitle.String
		}
		if areaID.Valid {
			t.AreaID = areaID.String
		}
		if areaTitle.Valid {
			t.AreaTitle = areaTitle.String
		}
		if headingID.Valid {
			t.HeadingID = headingID.String
		}
		if headingTitle.Valid {
			t.HeadingTitle = headingTitle.String
		}
		if tagTitles.Valid && tagTitles.String != "" {
			t.Tags = strings.Split(tagTitles.String, tagSeparator)
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

// escapeLike escapes SQL LIKE metacharacters so they match literally.
func escapeLike(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}
