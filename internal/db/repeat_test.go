package db

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"howett.net/plist"
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
	if start != 2 {
		t.Fatalf("expected start=2, got %d", start)
	}
	if !startDay.Valid || startDay.Int64 != 123 {
		t.Fatalf("expected startDate=123, got %v", startDay)
	}
	if bucket != 0 {
		t.Fatalf("expected startBucket=0, got %d", bucket)
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
	if err := store.ApplyRepeatRule("MISSING", update); err == nil {
		t.Fatal("expected applying to a missing UUID to fail")
	}
	if err := store.ClearRepeatRule("MISSING"); err == nil {
		t.Fatal("expected clearing a missing UUID to fail")
	}
}

func TestRepeatStateByIDDecodesCanonicalSemantics(t *testing.T) {
	path := filepath.Join(t.TempDir(), "things.sqlite3")
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	_, err = conn.Exec(`CREATE TABLE TMTask (
		uuid TEXT PRIMARY KEY, start INTEGER, startBucket INTEGER,
		deadline INTEGER,
		rt1_recurrenceRule BLOB, rt1_instanceCreationPaused INTEGER,
		rt1_nextInstanceStartDate INTEGER
	)`)
	if err != nil {
		t.Fatal(err)
	}
	anchor := time.Date(2025, 3, 4, 0, 0, 0, 0, time.Local)
	end := anchor.AddDate(0, 1, 0)
	rule, err := plist.Marshal(map[string]any{
		"tp": 0, "fu": 16, "fa": 2, "ia": float64(anchor.Unix()),
		"ed": float64(end.Unix()), "ts": -3,
	}, plist.XMLFormat)
	if err != nil {
		t.Fatal(err)
	}
	next := thingsDate(anchor.AddDate(0, 0, 2))
	_, err = conn.Exec(`INSERT INTO TMTask VALUES ('T1', 2, 0, 262213760, ?, 0, ?)`, rule, next)
	if err != nil {
		t.Fatal(err)
	}
	_, err = conn.Exec(`INSERT INTO TMTask VALUES ('BAD', 2, 0, NULL, X'0102', 0, NULL)`)
	if err != nil {
		t.Fatal(err)
	}
	normalDeadlineRule, err := plist.Marshal(map[string]any{
		"tp": 0, "fu": 16, "fa": 1, "ia": float64(anchor.Unix()), "ts": 0,
	}, plist.XMLFormat)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(`INSERT INTO TMTask VALUES ('NORMAL-DEADLINE', 2, 0, 132710400, ?, 0, NULL)`, normalDeadlineRule); err != nil {
		t.Fatal(err)
	}
	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
	store, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	state, err := store.RepeatStateByID("T1")
	if err != nil {
		t.Fatal(err)
	}
	if !state.Active || !state.Scheduled || state.Mode != "schedule" || state.Unit != "day" || state.Interval != 2 {
		t.Fatalf("unexpected semantic state: %#v", state)
	}
	if state.Anchor != "2025-03-04" || state.EndDate != "2025-04-04" || state.DeadlineOffset == nil || *state.DeadlineOffset != 3 {
		t.Fatalf("unexpected dates/offset: %#v", state)
	}
	bad, err := store.RepeatStateByID("BAD")
	if err != nil {
		t.Fatal(err)
	}
	if bad.DecodeWarning == "" {
		t.Fatalf("expected corrupt rule warning: %#v", bad)
	}
	normal, err := store.RepeatStateByID("NORMAL-DEADLINE")
	if err != nil {
		t.Fatal(err)
	}
	if normal.DeadlineOffset != nil {
		t.Fatalf("ordinary deadline decoded as repeat offset: %#v", normal)
	}
}

func TestDecodeRepeatRulePreservesExistingWarnings(t *testing.T) {
	rule, err := plist.Marshal(map[string]any{"tp": 99, "fu": 99}, plist.XMLFormat)
	if err != nil {
		t.Fatal(err)
	}
	state := &RepeatState{DecodeWarning: "repeat lifecycle state unavailable"}
	decodeRepeatRule(rule, false, state)
	if got, want := state.DecodeWarning, "repeat lifecycle state unavailable; unsupported recurrence mode; unsupported recurrence unit"; got != want {
		t.Fatalf("warning = %q, want %q", got, want)
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
