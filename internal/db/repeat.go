package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// RepeatTarget captures minimal metadata for repeat operations.
type RepeatTarget struct {
	UUID                string
	Title               string
	Type                int
	Status              int
	Trashed             bool
	Repeating           bool
	RepeatingTemplateID string
}

// RepeatUpdate describes the repeat fields to apply to a task or project.
type RepeatUpdate struct {
	RecurrenceRule            []byte
	InstanceCreationStartDate int
	InstanceCreationPaused    int
	InstanceCreationCount     int
	AfterCompletionReference  *int
	NextInstanceStartDate     *int
	Deadline                  *int
	SetDeadline               bool
}

// RepeatTargetByID returns repeat-related metadata for the given task ID.
func (s *Store) RepeatTargetByID(id string) (*RepeatTarget, error) {
	if s == nil || s.conn == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	if strings.TrimSpace(id) == "" {
		return nil, sql.ErrNoRows
	}
	var target RepeatTarget
	var repeating sql.NullInt64
	var repeatingTemplate sql.NullString
	if err := s.conn.QueryRow(
		`SELECT uuid, title, type, status, trashed, (rt1_recurrenceRule IS NOT NULL), rt1_repeatingTemplate
		 FROM TMTask WHERE uuid = ?`,
		id,
	).Scan(&target.UUID, &target.Title, &target.Type, &target.Status, &target.Trashed, &repeating, &repeatingTemplate); err != nil {
		return nil, err
	}
	if repeating.Valid {
		target.Repeating = repeating.Int64 != 0
	}
	if repeatingTemplate.Valid {
		target.RepeatingTemplateID = repeatingTemplate.String
	}
	return &target, nil
}

// TasksByTitleSince returns tasks created or modified with the given title after the timestamp.
func (s *Store) TasksByTitleSince(title string, taskType int, since float64) ([]TaskMatch, error) {
	if s == nil || s.conn == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	if strings.TrimSpace(title) == "" {
		return nil, fmt.Errorf("title required")
	}
	rows, err := s.conn.Query(
		`SELECT uuid,
			CASE
				WHEN userModificationDate IS NOT NULL AND userModificationDate > 0 THEN userModificationDate
				ELSE creationDate
			END AS created
		 FROM TMTask
		 WHERE type = ? AND lower(title) = lower(?) AND (creationDate >= ? OR userModificationDate >= ?)
		 ORDER BY created DESC`,
		taskType,
		title,
		since,
		since,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []TaskMatch
	for rows.Next() {
		var match TaskMatch
		if err := rows.Scan(&match.UUID, &match.Created); err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}
	return matches, rows.Err()
}

// TaskMatch identifies a task result for matching repeat updates after creation.
type TaskMatch struct {
	UUID    string
	Created float64
}

// ApplyRepeatRule updates a task/project with a recurrence rule in the database.
func (s *Store) ApplyRepeatRule(id string, update RepeatUpdate) error {
	if s == nil || s.conn == nil {
		return fmt.Errorf("database not initialized")
	}
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("task id required")
	}
	if len(update.RecurrenceRule) == 0 {
		return fmt.Errorf("recurrence rule required")
	}
	modified := float64(time.Now().Unix())

	var b strings.Builder
	b.WriteString("UPDATE TMTask SET ")
	b.WriteString("rt1_recurrenceRule = ?, ")
	b.WriteString("rt1_instanceCreationStartDate = ?, ")
	b.WriteString("rt1_instanceCreationPaused = ?, ")
	b.WriteString("rt1_instanceCreationCount = ?, ")
	b.WriteString("rt1_afterCompletionReferenceDate = ?, ")
	b.WriteString("rt1_nextInstanceStartDate = ?, ")
	if update.SetDeadline {
		b.WriteString("deadline = ?, ")
		b.WriteString("deadlineSuppressionDate = NULL, ")
	}
	b.WriteString("userModificationDate = ? ")
	b.WriteString("WHERE uuid = ?")

	args := []any{
		update.RecurrenceRule,
		update.InstanceCreationStartDate,
		update.InstanceCreationPaused,
		update.InstanceCreationCount,
		update.AfterCompletionReference,
		update.NextInstanceStartDate,
	}
	if update.SetDeadline {
		args = append(args, update.Deadline)
	}
	args = append(args, modified, id)

	_, err := s.conn.Exec(b.String(), args...)
	return err
}

// ClearRepeatRule removes recurrence metadata for a task/project.
func (s *Store) ClearRepeatRule(id string) error {
	if s == nil || s.conn == nil {
		return fmt.Errorf("database not initialized")
	}
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("task id required")
	}
	modified := float64(time.Now().Unix())
	_, err := s.conn.Exec(
		`UPDATE TMTask SET
			rt1_recurrenceRule = NULL,
			rt1_instanceCreationStartDate = NULL,
			rt1_instanceCreationPaused = 0,
			rt1_instanceCreationCount = 0,
			rt1_afterCompletionReferenceDate = NULL,
			rt1_nextInstanceStartDate = NULL,
			deadline = NULL,
			deadlineSuppressionDate = NULL,
			userModificationDate = ?
		  WHERE uuid = ?`,
		modified,
		id,
	)
	return err
}
