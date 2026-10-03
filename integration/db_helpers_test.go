package integration_test

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	_ "modernc.org/sqlite"
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
	yesterday := thingsDate(now.AddDate(0, 0, -1))
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
	if _, err := conn.Exec(`INSERT INTO TMTag (uuid, title) VALUES ('TAG_PLAN', 'PLAN');`); err != nil {
		t.Fatalf("insert plan tag: %v", err)
	}
	if _, err := conn.Exec(`INSERT INTO TMTaskTag (tasks, tags) VALUES ('T1', 'TAG1');`); err != nil {
		t.Fatalf("insert task tag: %v", err)
	}
	if _, err := conn.Exec(`INSERT INTO TMTaskTag (tasks, tags) VALUES ('P1', 'TAG_PLAN');`); err != nil {
		t.Fatalf("insert project tag: %v", err)
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
	if _, err := conn.Exec(`UPDATE TMTask SET startBucket = 0, todayIndex = 1, todayIndexReferenceDate = 135004288 WHERE uuid = 'TODAY1'`); err != nil {
		t.Fatalf("set Today ordering metadata: %v", err)
	}

	if _, err := conn.Exec(
		`INSERT INTO TMTask (uuid, type, status, trashed, title, area, start, creationDate, rt1_recurrenceRule) VALUES ('TPL1', ?, ?, 0, 'Template Task', 'A1', 1, ?, X'01')`,
		0,
		0,
		nowUnix,
	); err != nil {
		t.Fatalf("insert template task: %v", err)
	}
	if _, err := conn.Exec(
		`INSERT INTO TMTask (uuid, type, status, trashed, title, area, start, creationDate, rt1_repeatingTemplate) VALUES ('GEN1', ?, ?, 0, 'Generated Template Instance', 'A1', 1, ?, 'TPL1')`,
		0,
		0,
		nowUnix,
	); err != nil {
		t.Fatalf("insert generated repeating instance: %v", err)
	}

	// Projects the user closed, each still holding an open to-do. Things treats
	// those leftovers as archived and hides them from the active lists.
	closedProjects := []struct {
		uuid   string
		title  string
		status int
	}{
		{"PDONE", "Completed Project", 3},
		{"PCANC", "Canceled Project", 2},
	}
	for _, project := range closedProjects {
		if _, err := conn.Exec(
			`INSERT INTO TMTask (uuid, type, status, trashed, title, area, start, creationDate, stopDate) VALUES (?, ?, ?, 0, ?, 'A1', 1, ?, ?)`,
			project.uuid, 1, project.status, project.title, nowUnix, nowUnix,
		); err != nil {
			t.Fatalf("insert closed project %s: %v", project.uuid, err)
		}
	}

	closedProjectChildren := []struct {
		uuid      string
		title     string
		project   string
		status    int
		start     int
		startDate *int
		stopDate  *float64
	}{
		{"LEFTTODAY1", "Open Task in Completed Project", "PDONE", 0, 1, &today, nil},
		{"LEFTSCHED1", "Scheduled Task in Completed Project", "PDONE", 0, 2, &yesterday, nil},
		{"LEFTCANC1", "Open Task in Canceled Project", "PCANC", 0, 1, nil, nil},
		{"DONECHILD1", "Completed Task in Completed Project", "PDONE", 3, 1, nil, floatPtr(nowUnix)},
	}
	for _, child := range closedProjectChildren {
		if _, err := conn.Exec(
			`INSERT INTO TMTask (uuid, type, status, trashed, title, project, area, start, startDate, startBucket, creationDate, stopDate) VALUES (?, ?, ?, 0, ?, ?, 'A1', ?, ?, 0, ?, ?)`,
			child.uuid, 0, child.status, child.title, child.project, child.start, child.startDate, nowUnix, child.stopDate,
		); err != nil {
			t.Fatalf("insert closed project child %s: %v", child.uuid, err)
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
