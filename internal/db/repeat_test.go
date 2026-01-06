package db

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func TestApplyAndClearRepeatRule(t *testing.T) {
	path := filepath.Join(t.TempDir(), "things.sqlite3")
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if _, err := conn.Exec(`CREATE TABLE TMTask (
		uuid TEXT PRIMARY KEY,
		title TEXT,
		type INTEGER,
		status INTEGER,
		trashed INTEGER,
		start INTEGER,
		startDate INTEGER,
		startBucket INTEGER,
		deadline INTEGER,
		deadlineSuppressionDate INTEGER,
		rt1_repeatingTemplate TEXT,
		rt1_recurrenceRule BLOB,
		rt1_instanceCreationStartDate INTEGER,
		rt1_instanceCreationPaused INTEGER,
		rt1_instanceCreationCount INTEGER,
		rt1_afterCompletionReferenceDate INTEGER,
		rt1_nextInstanceStartDate INTEGER,
		userModificationDate REAL
	);`); err != nil {
		t.Fatalf("create schema: %v", err)
	}
	if _, err := conn.Exec(`INSERT INTO TMTask (uuid, title, type, status, trashed, start, startDate, startBucket) VALUES ('T1', 'Test', ?, ?, 0, 1, 123, 4);`, TaskTypeTodo, StatusIncomplete); err != nil {
		t.Fatalf("insert task: %v", err)
	}
	if err := conn.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}

	store, err := OpenWritable(path)
	if err != nil {
		t.Fatalf("open writable: %v", err)
	}
	defer store.Close()

	deadline := 262213760
	update := RepeatUpdate{
		RecurrenceRule:            []byte{0x01, 0x02},
		InstanceCreationStartDate: 123,
		InstanceCreationPaused:    0,
		InstanceCreationCount:     0,
		AfterCompletionReference:  nil,
		NextInstanceStartDate:     nil,
		Deadline:                  &deadline,
		SetDeadline:               true,
	}
	if err := store.ApplyRepeatRule("T1", update); err != nil {
		t.Fatalf("apply repeat: %v", err)
	}

	var (
		start      int
		startDay   sql.NullInt64
		bucket     int
		rule       []byte
		dbDeadline sql.NullInt64
		modified   sql.NullFloat64
	)
	if err := store.conn.QueryRow(`SELECT start, startDate, startBucket, rt1_recurrenceRule, deadline, userModificationDate FROM TMTask WHERE uuid = 'T1'`).Scan(&start, &startDay, &bucket, &rule, &dbDeadline, &modified); err != nil {
		t.Fatalf("select updated: %v", err)
	}
	if start != 1 {
		t.Fatalf("expected start=1, got %d", start)
	}
	if !startDay.Valid || startDay.Int64 != 123 {
		t.Fatalf("expected startDate=123, got %v", startDay)
	}
	if bucket != 4 {
		t.Fatalf("expected startBucket=4, got %d", bucket)
	}
	if len(rule) == 0 {
		t.Fatalf("expected recurrence rule bytes")
	}
	if !dbDeadline.Valid || int(dbDeadline.Int64) != deadline {
		t.Fatalf("expected deadline %d, got %v", deadline, dbDeadline)
	}
	if !modified.Valid || modified.Float64 <= 0 {
		t.Fatalf("expected userModificationDate to be set")
	}

	if err := store.ClearRepeatRule("T1"); err != nil {
		t.Fatalf("clear repeat: %v", err)
	}
	var clearedRule []byte
	if err := store.conn.QueryRow(`SELECT rt1_recurrenceRule, deadline FROM TMTask WHERE uuid = 'T1'`).Scan(&clearedRule, &dbDeadline); err != nil {
		t.Fatalf("select cleared: %v", err)
	}
	if len(clearedRule) != 0 {
		t.Fatalf("expected recurrence rule cleared")
	}
	if dbDeadline.Valid {
		t.Fatalf("expected deadline cleared")
	}
}

func TestTasksByTitleSinceUsesModificationDate(t *testing.T) {
	path := filepath.Join(t.TempDir(), "things.sqlite3")
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if _, err := conn.Exec(`CREATE TABLE TMTask (
		uuid TEXT PRIMARY KEY,
		title TEXT,
		type INTEGER,
		creationDate REAL,
		userModificationDate REAL
	);`); err != nil {
		t.Fatalf("create schema: %v", err)
	}

	now := float64(time.Now().Unix())
	old := now - 120
	if _, err := conn.Exec(`INSERT INTO TMTask (uuid, title, type, creationDate, userModificationDate) VALUES ('T1', 'Repeat Me', ?, ?, ?);`, TaskTypeTodo, old, now); err != nil {
		t.Fatalf("insert task: %v", err)
	}
	if _, err := conn.Exec(`INSERT INTO TMTask (uuid, title, type, creationDate, userModificationDate) VALUES ('T2', 'Repeat Me', ?, ?, ?);`, TaskTypeTodo, old, old); err != nil {
		t.Fatalf("insert task: %v", err)
	}
	if err := conn.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}

	store, err := Open(path)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer store.Close()

	matches, err := store.TasksByTitleSince("Repeat Me", TaskTypeTodo, now-10)
	if err != nil {
		t.Fatalf("TasksByTitleSince: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0].UUID != "T1" {
		t.Fatalf("expected T1, got %s", matches[0].UUID)
	}
}
