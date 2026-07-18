package cli

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"
)

func writeTestDB(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "Things.sqlite3")
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer conn.Close()

	statements := []string{
		`CREATE TABLE TMArea (uuid TEXT PRIMARY KEY, title TEXT, visible INTEGER, "index" INTEGER);`,
		`CREATE TABLE TMTask (
			uuid TEXT PRIMARY KEY,
			type INTEGER,
			status INTEGER,
			trashed INTEGER,
			title TEXT,
			notes TEXT,
			area TEXT,
			project TEXT,
			heading TEXT,
			start INTEGER,
			startDate INTEGER,
			startBucket INTEGER,
			deadline INTEGER,
			deadlineSuppressionDate INTEGER,
			creationDate REAL,
			userModificationDate REAL,
			stopDate REAL,
			"index" INTEGER,
			rt1_repeatingTemplate TEXT,
			rt1_recurrenceRule BLOB,
			rt1_instanceCreationStartDate INTEGER,
			rt1_instanceCreationPaused INTEGER,
			rt1_instanceCreationCount INTEGER,
			rt1_afterCompletionReferenceDate INTEGER,
			rt1_nextInstanceStartDate INTEGER,
			repeater BLOB,
			todayIndex INTEGER,
			todayIndexReferenceDate INTEGER
		);`,
		`CREATE TABLE TMTag (uuid TEXT PRIMARY KEY, title TEXT, shortcut TEXT, parent TEXT);`,
		`CREATE TABLE TMTaskTag (tasks TEXT NOT NULL, tags TEXT NOT NULL);`,
		`CREATE TABLE TMChecklistItem (
			uuid TEXT PRIMARY KEY,
			userModificationDate REAL,
			creationDate REAL,
			title TEXT,
			status INTEGER,
			stopDate REAL,
			"index" INTEGER,
			task TEXT
		);`,
	}
	for _, stmt := range statements {
		if _, err := conn.Exec(stmt); err != nil {
			t.Fatalf("create schema: %v", err)
		}
	}

	if _, err := conn.Exec(`INSERT INTO TMArea (uuid, title, visible, "index") VALUES ('A1', 'Home', 1, 1);`); err != nil {
		t.Fatalf("insert area: %v", err)
	}
	if _, err := conn.Exec(`INSERT INTO TMTask (uuid, type, status, trashed, title, area) VALUES ('P1', ?, ?, 0, 'Project One', 'A1');`, 1, 0); err != nil {
		t.Fatalf("insert project: %v", err)
	}
	now := time.Now()
	today := thingsDate(now)
	tomorrow := thingsDate(now.AddDate(0, 0, 1))
	nowUnix := float64(now.Unix())

	if _, err := conn.Exec(`INSERT INTO TMTask (uuid, type, status, trashed, title, project, area, heading, notes, start, creationDate) VALUES ('T1', ?, ?, 0, 'Task One', 'P1', 'A1', 'H1', 'Some notes', 1, ?);`, 0, 0, nowUnix); err != nil {
		t.Fatalf("insert task: %v", err)
	}
	if _, err := conn.Exec(`INSERT INTO TMTask (uuid, type, status, trashed, title, project, area, heading, notes) VALUES ('H1', ?, ?, 0, 'Heading', 'P1', 'A1', '', '');`, 2, 0); err != nil {
		t.Fatalf("insert heading: %v", err)
	}
	if _, err := conn.Exec(`INSERT INTO TMTag (uuid, title) VALUES ('TAG1', 'urgent');`); err != nil {
		t.Fatalf("insert tag: %v", err)
	}
	if _, err := conn.Exec(`INSERT INTO TMTaskTag (tasks, tags) VALUES ('T1', 'TAG1');`); err != nil {
		t.Fatalf("insert task tag: %v", err)
	}
	if _, err := conn.Exec(`INSERT INTO TMChecklistItem (uuid, title, status, "index", task) VALUES ('C1', 'Check Item', 0, 0, 'T1');`); err != nil {
		t.Fatalf("insert checklist: %v", err)
	}

	inserts := []struct {
		uuid      string
		title     string
		status    int
		trashed   int
		start     int
		startDate *int
		deadline  *int
		stopDate  *float64
	}{
		{"INBOX1", "Inbox Task", 0, 0, 0, nil, nil, nil},
		{"ANY1", "Anytime Task", 0, 0, 1, nil, nil, nil},
		{"TODAY1", "Today Task", 0, 0, 1, &today, nil, nil},
		{"UP1", "Upcoming Task", 0, 0, 2, &tomorrow, nil, nil},
		{"SOM1", "Someday Task", 0, 0, 2, nil, nil, nil},
		{"DL1", "Deadline Task", 0, 0, 1, nil, &tomorrow, nil},
		{"COMP1", "Completed Task", 3, 0, 1, nil, nil, floatPtr(nowUnix)},
		{"CANC1", "Canceled Task", 2, 0, 1, nil, nil, floatPtr(nowUnix)},
		{"TRASH1", "Trashed Task", 0, 1, 1, nil, nil, nil},
	}
	for _, item := range inserts {
		if _, err := conn.Exec(
			`INSERT INTO TMTask (uuid, type, status, trashed, title, start, startDate, deadline, creationDate, stopDate) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			item.uuid,
			0,
			item.status,
			item.trashed,
			item.title,
			item.start,
			item.startDate,
			item.deadline,
			nowUnix,
			item.stopDate,
		); err != nil {
			t.Fatalf("insert task %s: %v", item.uuid, err)
		}
	}

	return path
}

func thingsDate(t time.Time) int {
	date := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return date.Year()<<16 | int(date.Month())<<12 | date.Day()<<7
}

func floatPtr(v float64) *float64 {
	return &v
}
