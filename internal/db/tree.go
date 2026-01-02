package db

import (
	"fmt"
	"strings"
)

// AreasTree returns a hierarchical area -> project -> heading -> todo tree.
func (s *Store) AreasTree(filter TaskFilter, onlyProjects bool) ([]TreeItem, error) {
	if s == nil || s.conn == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	areas, err := s.Areas()
	if err != nil {
		return nil, err
	}
	items := make([]TreeItem, 0, len(areas))
	for _, area := range areas {
		if filter.AreaID != "" && filter.AreaID != area.UUID {
			continue
		}
		areaItem := TreeItem{
			UUID:  area.UUID,
			Type:  "area",
			Title: area.Title,
		}

		projectFilter := filter
		projectFilter.AreaID = area.UUID
		projects, err := s.queryTaskItems(TaskTypeProject, "", nil, projectFilter, "t.title COLLATE NOCASE")
		if err != nil {
			return nil, err
		}

		for _, project := range projects {
			projectItem := project
			if !onlyProjects {
				children, err := s.projectChildren(project.UUID, filter)
				if err != nil {
					return nil, err
				}
				projectItem.Items = children
			}
			areaItem.Items = append(areaItem.Items, projectItem)
		}

		if !onlyProjects {
			taskFilter := filter
			taskFilter.AreaID = area.UUID
			taskFilter.Types = []int{TaskTypeTodo}
			tasks, err := s.queryTasks("t.area = ? AND t.project IS NULL", []any{area.UUID}, taskFilter, "")
			if err != nil {
				return nil, err
			}
			for _, task := range tasks {
				areaItem.Items = append(areaItem.Items, taskToTree(task))
			}
		}

		items = append(items, areaItem)
	}
	return items, nil
}

// ProjectsTree returns a hierarchical project -> heading -> todo tree.
func (s *Store) ProjectsTree(filter TaskFilter, onlyProjects bool) ([]TreeItem, error) {
	if s == nil || s.conn == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	projects, err := s.queryTaskItems(TaskTypeProject, "", nil, filter, "t.title COLLATE NOCASE")
	if err != nil {
		return nil, err
	}

	for i := range projects {
		if onlyProjects {
			continue
		}
		children, err := s.projectChildren(projects[i].UUID, filter)
		if err != nil {
			return nil, err
		}
		projects[i].Items = children
	}

	return projects, nil
}

// ProjectsWithoutAreaTree returns a hierarchical project -> heading -> todo tree for projects with no area.
func (s *Store) ProjectsWithoutAreaTree(filter TaskFilter, onlyProjects bool) ([]TreeItem, error) {
	if s == nil || s.conn == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	projects, err := s.queryTaskItems(TaskTypeProject, "t.area IS NULL", nil, filter, "t.title COLLATE NOCASE")
	if err != nil {
		return nil, err
	}

	for i := range projects {
		if onlyProjects {
			continue
		}
		children, err := s.projectChildren(projects[i].UUID, filter)
		if err != nil {
			return nil, err
		}
		projects[i].Items = children
	}

	return projects, nil
}

func (s *Store) projectChildren(projectID string, filter TaskFilter) ([]TreeItem, error) {
	children := make([]TreeItem, 0, 16)
	headingFilter := filter
	headingFilter.ProjectID = projectID
	headings, err := s.queryTaskItems(TaskTypeHeading, "", nil, headingFilter, "t.title COLLATE NOCASE")
	if err != nil {
		return nil, err
	}

	for _, heading := range headings {
		taskFilter := filter
		taskFilter.ProjectID = projectID
		taskFilter.Types = []int{TaskTypeTodo}
		tasks, err := s.queryTasks("t.project = ? AND t.heading = ?", []any{projectID, heading.UUID}, taskFilter, "")
		if err != nil {
			return nil, err
		}
		for _, task := range tasks {
			heading.Items = append(heading.Items, taskToTree(task))
		}
		if len(heading.Items) > 0 {
			children = append(children, heading)
		}
	}

	taskFilter := filter
	taskFilter.ProjectID = projectID
	taskFilter.Types = []int{TaskTypeTodo}
	tasks, err := s.queryTasks("t.project = ? AND t.heading IS NULL", []any{projectID}, taskFilter, "")
	if err != nil {
		return nil, err
	}
	for _, task := range tasks {
		children = append(children, taskToTree(task))
	}

	return children, nil
}

func (s *Store) queryTaskItems(taskType int, where string, args []any, filter TaskFilter, order string) ([]TreeItem, error) {
	var b strings.Builder
	b.WriteString("SELECT t.uuid, t.title, t.status, t.trashed ")
	b.WriteString("FROM TMTask t ")
	b.WriteString("LEFT JOIN TMTask p ON t.project = p.uuid ")
	b.WriteString("WHERE t.type = ?")

	params := []any{taskType}
	if where != "" {
		b.WriteString(" AND (" + where + ")")
		params = append(params, args...)
	}

	if !filter.IncludeTrashed {
		b.WriteString(" AND t.trashed = 0")
	}
	if filter.ExcludeTrashedContext {
		b.WriteString(" AND NOT IFNULL(p.trashed, 0)")
	}
	if filter.Status != nil {
		b.WriteString(" AND t.status = ?")
		params = append(params, *filter.Status)
	}
	if filter.AreaID != "" {
		b.WriteString(" AND t.area = ?")
		params = append(params, filter.AreaID)
	}
	if filter.ProjectID != "" {
		if taskType == TaskTypeProject {
			b.WriteString(" AND t.uuid = ?")
		} else {
			b.WriteString(" AND t.project = ?")
		}
		params = append(params, filter.ProjectID)
	}
	if filter.TagID != "" {
		b.WriteString(" AND EXISTS (SELECT 1 FROM TMTaskTag tt WHERE tt.tasks = t.uuid AND tt.tags = ?)")
		params = append(params, filter.TagID)
	}
	b.WriteString(" AND t.rt1_recurrenceRule IS NULL")

	if order == "" {
		order = "t.\"index\""
	}
	b.WriteString(" ORDER BY " + order)

	rows, err := s.conn.Query(b.String(), params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]TreeItem, 0, 32)
	for rows.Next() {
		var item TreeItem
		var status int
		var trashed int
		if err := rows.Scan(&item.UUID, &item.Title, &status, &trashed); err != nil {
			return nil, err
		}
		item.Type = taskTypeLabel(taskType)
		item.Status = &status
		trashedBool := trashed != 0
		item.Trashed = &trashedBool
		items = append(items, item)
	}
	return items, rows.Err()
}

func taskToTree(task Task) TreeItem {
	status := task.Status
	trashed := task.Trashed
	return TreeItem{
		UUID:    task.UUID,
		Type:    "to-do",
		Title:   task.Title,
		Status:  &status,
		Trashed: &trashed,
	}
}
