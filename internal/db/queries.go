package db

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"
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

// TaskByID returns a single task by UUID.
func (s *Store) TaskByID(id string) (*Task, error) {
	if strings.TrimSpace(id) == "" {
		return nil, sql.ErrNoRows
	}
	filter := TaskFilter{
		IncludeTrashed:        true,
		ExcludeTrashedContext: false,
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
	regular, err := s.queryTasks("t.start = 1 AND t.startDate IS NOT NULL", nil, filter, "t.todayIndex")
	if err != nil {
		return nil, err
	}
	unconfirmedScheduled, err := s.queryTasks("t.start = 2 AND t.startDate IS NOT NULL AND t.startDate <= "+todayExpr, nil, filter, "t.todayIndex")
	if err != nil {
		return nil, err
	}
	unconfirmedOverdue, err := s.queryTasks("t.startDate IS NULL AND t.deadline IS NOT NULL AND t.deadline <= "+todayExpr+" AND t.deadlineSuppressionDate IS NULL", nil, filter, "t.todayIndex")
	if err != nil {
		return nil, err
	}
	result := append(regular, unconfirmedScheduled...)
	result = append(result, unconfirmedOverdue...)
	maxIndex := int(^uint(0) >> 1)
	sort.Slice(result, func(i, j int) bool {
		left := maxIndex
		right := maxIndex
		if result[i].TodayIndex != nil {
			left = *result[i].TodayIndex
		}
		if result[j].TodayIndex != nil {
			right = *result[j].TodayIndex
		}
		if left != right {
			return left < right
		}
		return result[i].StartDate < result[j].StartDate
	})
	return result, nil
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
func (s *Store) ItemByID(id string) (*Item, error) {
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
func (s *Store) ItemsByTitle(title string) ([]Item, error) {
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
	return "", fmt.Errorf("area not found: %s", input)
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
	return "", fmt.Errorf("project not found: %s", input)
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
	return "", fmt.Errorf("tag not found: %s", input)
}

func thingsDateTodayExpr() string {
	return "((strftime('%Y', date('now', 'localtime')) << 16) | (strftime('%m', date('now', 'localtime')) << 12) | (strftime('%d', date('now', 'localtime')) << 7))"
}

func (s *Store) queryTasks(where string, args []any, filter TaskFilter, order string) ([]Task, error) {
	if s == nil || s.conn == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	var b strings.Builder
	b.WriteString("SELECT t.uuid, t.type, t.title, t.status, t.trashed, t.notes, t.start, t.startDate, t.deadline, t.stopDate, t.creationDate, t.userModificationDate, t.\"index\", t.todayIndex, ")
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
		b.WriteString(" AND (t.project = ? OR hp.uuid = ?)")
		params = append(params, filter.ProjectID, filter.ProjectID)
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
		b.WriteString(" AND (lower(t.title) LIKE lower(?) OR lower(t.notes) LIKE lower(?) OR lower(a.title) LIKE lower(?))")
		like := "%" + filter.Search + "%"
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
	b.WriteString(" AND t.rt1_recurrenceRule IS NULL")

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
		var deadline sql.NullInt64
		var stopDate sql.NullFloat64
		var created sql.NullFloat64
		var modified sql.NullFloat64
		var index sql.NullInt64
		var todayIndex sql.NullInt64
		var projectID sql.NullString
		var projectTitle sql.NullString
		var areaID sql.NullString
		var areaTitle sql.NullString
		var headingID sql.NullString
		var headingTitle sql.NullString
		var tagTitles sql.NullString
		if err := rows.Scan(&t.UUID, &taskType, &t.Title, &t.Status, &t.Trashed, &notes, &start, &startDate, &deadline, &stopDate, &created, &modified, &index, &todayIndex, &projectID, &projectTitle, &areaID, &areaTitle, &headingID, &headingTitle, &tagTitles); err != nil {
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
		if notes.Valid {
			t.Notes = notes.String
		}
		if start.Valid {
			t.Start = startLabel(int(start.Int64))
		}
		if startDate.Valid {
			t.StartDate = formatThingsDate(startDate.Int64)
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
