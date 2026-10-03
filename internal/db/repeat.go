package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"howett.net/plist"
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
	InstanceCount       int
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

// ValidateRepeatSchema checks the direct-write capabilities before any mutation.
func (s *Store) ValidateRepeatSchema() error {
	if s == nil || s.conn == nil {
		return fmt.Errorf("database not initialized")
	}
	rows, err := s.conn.Query(`PRAGMA table_info(TMTask)`)
	if err != nil {
		return err
	}
	defer rows.Close()
	found := map[string]bool{}
	for rows.Next() {
		var cid int
		var name, typ string
		var notNull, pk int
		var defaultValue any
		if err := rows.Scan(&cid, &name, &typ, &notNull, &defaultValue, &pk); err != nil {
			return err
		}
		found[name] = true
	}
	required := []string{"uuid", "start", "startBucket", "deadline", "deadlineSuppressionDate", "rt1_recurrenceRule", "rt1_instanceCreationStartDate", "rt1_instanceCreationPaused", "rt1_instanceCreationCount", "rt1_afterCompletionReferenceDate", "rt1_nextInstanceStartDate", "userModificationDate"}
	for _, name := range required {
		if !found[name] {
			return fmt.Errorf("unsupported Things database schema: missing TMTask.%s", name)
		}
	}
	return rows.Err()
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
	var instanceCount sql.NullInt64
	if err := s.conn.QueryRow(
		`SELECT uuid, title, type, status, trashed, (rt1_recurrenceRule IS NOT NULL), rt1_repeatingTemplate, COALESCE(rt1_instanceCreationCount, 0)
		 FROM TMTask WHERE uuid = ?`,
		id,
	).Scan(&target.UUID, &target.Title, &target.Type, &target.Status, &target.Trashed, &repeating, &repeatingTemplate, &instanceCount); err != nil {
		return nil, err
	}
	if repeating.Valid {
		target.Repeating = repeating.Int64 != 0
	}
	if repeatingTemplate.Valid {
		target.RepeatingTemplateID = repeatingTemplate.String
	}
	target.InstanceCount = int(instanceCount.Int64)
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
	b.WriteString("start = 2, ")
	b.WriteString("startBucket = 0, ")
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

	result, err := s.conn.Exec(b.String(), args...)
	if err != nil {
		return err
	}
	return requireAffectedRow(result, id)
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
	result, err := s.conn.Exec(
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
	if err != nil {
		return err
	}
	return requireAffectedRow(result, id)
}

func requireAffectedRow(result sql.Result, id string) error {
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("inspect repeat update for %s: %w", id, err)
	}
	if count != 1 {
		return fmt.Errorf("repeat update affected %d rows for %s", count, id)
	}
	return nil
}

// RepeatStateByID reads and decodes the verification-relevant state of a template.
// Unsupported rules remain inspectable and carry DecodeWarning.
func (s *Store) RepeatStateByID(id string) (*RepeatState, error) {
	if s == nil || s.conn == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	var rule []byte
	var start, bucket, paused, deadline sql.NullInt64
	var next sql.NullInt64
	err := s.conn.QueryRow(`SELECT rt1_recurrenceRule, start, startBucket, deadline
		FROM TMTask WHERE uuid = ?`, id).Scan(&rule, &start, &bucket, &deadline)
	if err != nil {
		return nil, err
	}
	// Older read-only fixtures and Things schemas may not expose the rt1 lifecycle
	// columns. Keep the template listable, but surface the missing capability.
	lifecycleErr := s.conn.QueryRow(`SELECT rt1_instanceCreationPaused,
		rt1_nextInstanceStartDate FROM TMTask WHERE uuid = ?`, id).Scan(&paused, &next)
	state := &RepeatState{Paused: paused.Valid && paused.Int64 != 0}
	if lifecycleErr != nil {
		appendRepeatWarning(state, fmt.Sprintf("repeat lifecycle state unavailable: %v", lifecycleErr))
	}
	if bucket.Valid {
		v := int(bucket.Int64)
		state.StartBucket = &v
	}
	state.Scheduled = start.Valid && start.Int64 == 2 && bucket.Valid && bucket.Int64 == 0
	state.Active = len(rule) > 0 && state.Scheduled && !state.Paused
	if next.Valid {
		state.NextDate = formatThingsDate(next.Int64)
	}
	if len(rule) == 0 {
		return state, nil
	}
	// Repeat deadlines use Things' year-4001 sentinel. An ordinary task
	// deadline must not make a rule's default ts=0 look like an explicit
	// repeat deadline offset.
	hasRepeatDeadline := deadline.Valid && int(deadline.Int64)>>16 == 4001
	decodeRepeatRule(rule, hasRepeatDeadline, state)
	return state, nil
}

// RepeatStatesByIDs hydrates repeat projections in one query. It is used by
// template listings to avoid one database round trip per task.
func (s *Store) RepeatStatesByIDs(ids []string) (map[string]*RepeatState, error) {
	states := make(map[string]*RepeatState, len(ids))
	if len(ids) == 0 {
		return states, nil
	}
	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")
	args := make([]any, len(ids))
	for i := range ids {
		args[i] = ids[i]
	}
	hasLifecycle := true
	var ignored int
	if err := s.conn.QueryRow(`SELECT rt1_instanceCreationPaused FROM TMTask LIMIT 1`).Scan(&ignored); err != nil && strings.Contains(err.Error(), "no such column") {
		hasLifecycle = false
	}
	pausedExpr, nextExpr := "rt1_instanceCreationPaused", "rt1_nextInstanceStartDate"
	if !hasLifecycle {
		pausedExpr, nextExpr = "NULL", "NULL"
	}
	rows, err := s.conn.Query(`SELECT uuid, rt1_recurrenceRule, start, startBucket,
		deadline, `+pausedExpr+`, `+nextExpr+`
		FROM TMTask WHERE uuid IN (`+placeholders+`)`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		var rule []byte
		var start, bucket, deadline, paused, next sql.NullInt64
		if err := rows.Scan(&id, &rule, &start, &bucket, &deadline, &paused, &next); err != nil {
			return nil, err
		}
		state := &RepeatState{Paused: paused.Valid && paused.Int64 != 0}
		if !hasLifecycle {
			appendRepeatWarning(state, "repeat lifecycle state unavailable")
		}
		if bucket.Valid {
			v := int(bucket.Int64)
			state.StartBucket = &v
		}
		state.Scheduled = start.Valid && start.Int64 == 2 && bucket.Valid && bucket.Int64 == 0
		state.Active = len(rule) > 0 && state.Scheduled && !state.Paused
		if next.Valid {
			state.NextDate = formatThingsDate(next.Int64)
		}
		if len(rule) > 0 {
			decodeRepeatRule(rule, deadline.Valid && int(deadline.Int64)>>16 == 4001, state)
		}
		states[id] = state
	}
	return states, rows.Err()
}

func decodeRepeatRule(rule []byte, hasDeadline bool, state *RepeatState) {
	var values map[string]any
	if _, err := plist.Unmarshal(rule, &values); err != nil {
		appendRepeatWarning(state, fmt.Sprintf("cannot decode recurrence rule: %v", err))
		return
	}
	mode, ok := plistInt(values["tp"])
	if !ok || (mode != 0 && mode != 1) {
		appendRepeatWarning(state, "unsupported recurrence mode")
	} else if mode == 0 {
		state.Mode = "schedule"
	} else {
		state.Mode = "after-completion"
	}
	unit, ok := plistInt(values["fu"])
	units := map[int]string{16: "day", 256: "week", 8: "month", 4: "year"}
	if name, exists := units[unit]; ok && exists {
		state.Unit = name
	} else {
		appendRepeatWarning(state, "unsupported recurrence unit")
	}
	if interval, ok := plistInt(values["fa"]); ok {
		state.Interval = interval
	}
	if anchor, ok := plistFloat(values["ia"]); ok {
		state.Anchor = time.Unix(int64(anchor), 0).In(time.Local).Format("2006-01-02")
	}
	if end, ok := plistFloat(values["ed"]); ok {
		t := time.Unix(int64(end), 0).In(time.Local)
		if t.Year() < 4000 {
			state.EndDate = t.Format("2006-01-02")
		}
	}
	if offset, ok := plistInt(values["ts"]); ok && hasDeadline {
		v := -offset
		state.DeadlineOffset = &v
		// The encoded interval anchor and end boundary are stored in the
		// shifted deadline frame; restore both to user-facing coordinates
		// (the dates supplied via --repeat-start / --repeat-until) for
		// verification and output. This matches Things' own encoding: the app
		// stores ia = requested anchor + offset so occurrences land on the
		// anchor weekday. Rules written by the released CLI set ia directly
		// with ts<0 and therefore placed occurrences offset days off the
		// anchor weekday; no such correct legacy encoding exists to preserve.
		if state.Anchor != "" {
			parsed, err := time.ParseInLocation("2006-01-02", state.Anchor, time.Local)
			if err == nil {
				state.Anchor = parsed.AddDate(0, 0, offset).Format("2006-01-02")
			}
		}
		if state.EndDate != "" {
			parsed, err := time.ParseInLocation("2006-01-02", state.EndDate, time.Local)
			if err == nil {
				state.EndDate = parsed.AddDate(0, 0, offset).Format("2006-01-02")
			}
		}
	}
	if count, ok := plistInt(values["rc"]); ok && count > 0 {
		v := count
		state.Count = &v
	}
}

func appendRepeatWarning(state *RepeatState, warning string) {
	if state.DecodeWarning == "" {
		state.DecodeWarning = warning
		return
	}
	state.DecodeWarning += "; " + warning
}

func plistInt(value any) (int, bool) {
	switch v := value.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case uint64:
		return int(v), true
	case float64:
		return int(v), true
	default:
		return 0, false
	}
}

func plistFloat(value any) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case int64:
		return float64(v), true
	case uint64:
		return float64(v), true
	default:
		return 0, false
	}
}
