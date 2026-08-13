package db

import (
	"database/sql"
	"strings"
	"testing"
	"time"
)

func TestProjectsAndTasksQueries(t *testing.T) {
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer conn.Close()

	if err := seedTestDB(conn); err != nil {
		t.Fatalf("seed db: %v", err)
	}

	store := &Store{conn: conn, path: ":memory:"}
	status := StatusIncomplete
	projects, err := store.Projects(ProjectFilter{Status: &status})
	if err != nil {
		t.Fatalf("projects: %v", err)
	}
	if len(projects) != 1 || projects[0].Title != "Project One" {
		t.Fatalf("unexpected projects: %#v", projects)
	}

	projectID, err := store.ResolveProjectID("Project One")
	if err != nil {
		t.Fatalf("resolve project: %v", err)
	}
	if projectID == "" {
		t.Fatalf("expected project id")
	}

	tags, err := store.Tags()
	if err != nil {
		t.Fatalf("tags: %v", err)
	}
	if len(tags) != 1 || tags[0].Title != "urgent" || tags[0].Usage != 1 {
		t.Fatalf("unexpected tags: %#v", tags)
	}

	tasks, err := store.Tasks(TaskFilter{ProjectID: projectID, Status: &status, ExcludeTrashedContext: true, Types: []int{TaskTypeTodo}})
	if err != nil {
		t.Fatalf("tasks: %v", err)
	}
	if len(tasks) != 1 || tasks[0].Title != "Task One" {
		t.Fatalf("unexpected tasks: %#v", tasks)
	}
	if tasks[0].TodayIndexReferenceDate == nil || *tasks[0].TodayIndexReferenceDate != 135004288 {
		t.Fatalf("unexpected today index reference date: %#v", tasks[0].TodayIndexReferenceDate)
	}
	projectTask, err := store.TaskByID("P1")
	if err != nil {
		t.Fatalf("project task by ID: %v", err)
	}
	if projectTask.TodayIndexReferenceDate != nil {
		t.Fatalf("expected NULL reference date to remain nil, got %#v", projectTask.TodayIndexReferenceDate)
	}
	if tasks[0].Notes != "Some notes" {
		t.Fatalf("unexpected notes: %q", tasks[0].Notes)
	}
	if tasks[0].Start != "Anytime" {
		t.Fatalf("unexpected start: %q", tasks[0].Start)
	}
	if tasks[0].StartDate != "2025-01-02" {
		t.Fatalf("unexpected start date: %q", tasks[0].StartDate)
	}
	if tasks[0].Deadline != "2025-01-04" {
		t.Fatalf("unexpected deadline: %q", tasks[0].Deadline)
	}
	if tasks[0].Created != "2025-01-02 03:04:05" {
		t.Fatalf("unexpected created: %q", tasks[0].Created)
	}
	if tasks[0].Modified != "2025-01-02 03:04:05" {
		t.Fatalf("unexpected modified: %q", tasks[0].Modified)
	}
	if len(tasks[0].Tags) != 1 || tasks[0].Tags[0] != "urgent" {
		t.Fatalf("unexpected tags: %#v", tasks[0].Tags)
	}

	repeating, err := store.Tasks(TaskFilter{Status: &status, ExcludeTrashedContext: true, Types: []int{TaskTypeTodo}, RepeatingOnly: true})
	if err != nil {
		t.Fatalf("repeating tasks: %v", err)
	}
	if len(repeating) != 1 || repeating[0].Title != "Repeat Task" || !repeating[0].Repeating {
		t.Fatalf("unexpected repeating tasks: %#v", repeating)
	}

	templates, err := store.TemplatesTasks(TaskFilter{Status: &status, ExcludeTrashedContext: true, Types: []int{TaskTypeTodo}})
	if err != nil {
		t.Fatalf("template tasks: %v", err)
	}
	if len(templates) != 1 || templates[0].Title != "Repeat Task" || !templates[0].Repeating {
		t.Fatalf("unexpected template tasks: %#v", templates)
	}

	if _, err := conn.Exec(`UPDATE TMTask SET rt1_recurrenceRule = X'01', start = 2, startBucket = 0 WHERE uuid = 'P1'`); err != nil {
		t.Fatalf("make project repeating: %v", err)
	}
	templatesWithProjects, err := store.TemplatesTasks(TaskFilter{Status: &status, ExcludeTrashedContext: true, Types: []int{TaskTypeTodo, TaskTypeProject}})
	if err != nil {
		t.Fatalf("template tasks with projects: %v", err)
	}
	if len(templatesWithProjects) != 2 {
		t.Fatalf("expected todo + project templates, got %#v", templatesWithProjects)
	}
	repeatingProjects, err := store.Tasks(TaskFilter{Status: &status, ExcludeTrashedContext: true, Types: []int{TaskTypeProject}, RepeatingOnly: true})
	if err != nil {
		t.Fatalf("repeating projects: %v", err)
	}
	if len(repeatingProjects) != 1 || repeatingProjects[0].Title != "Project One" || !repeatingProjects[0].Repeating {
		t.Fatalf("unexpected repeating projects: %#v", repeatingProjects)
	}

	searched, err := store.Tasks(TaskFilter{Search: "notes", Status: &status, ExcludeTrashedContext: true, Types: []int{TaskTypeTodo}})
	if err != nil {
		t.Fatalf("search tasks: %v", err)
	}
	if len(searched) != 1 || searched[0].Title != "Task One" {
		t.Fatalf("unexpected search: %#v", searched)
	}

	searchedArea, err := store.Tasks(TaskFilter{Search: "Home", Status: &status, ExcludeTrashedContext: true, Types: []int{TaskTypeTodo}})
	if err != nil {
		t.Fatalf("search tasks by area: %v", err)
	}
	if len(searchedArea) != 1 || searchedArea[0].Title != "Task One" {
		t.Fatalf("unexpected search by area: %#v", searchedArea)
	}

	withChecklist, err := store.Tasks(TaskFilter{ProjectID: projectID, Status: &status, IncludeChecklist: true, ExcludeTrashedContext: true, Types: []int{TaskTypeTodo}})
	if err != nil {
		t.Fatalf("tasks with checklist: %v", err)
	}
	if len(withChecklist) != 1 || len(withChecklist[0].Checklist) != 1 || withChecklist[0].Checklist[0].Title != "Check Item" {
		t.Fatalf("unexpected checklist: %#v", withChecklist)
	}
}

func seedTestDB(conn *sql.DB) error {
	now := time.Date(2025, 1, 2, 3, 4, 5, 0, time.Local)
	startDate := thingsDateForTest(now)
	deadline := thingsDateForTest(now.AddDate(0, 0, 2))
	nowUnix := float64(now.Unix())

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
			return err
		}
	}

	if _, err := conn.Exec(`INSERT INTO TMArea (uuid, title, visible, "index") VALUES ('A1', 'Home', 1, 1);`); err != nil {
		return err
	}
	if _, err := conn.Exec(`INSERT INTO TMTask (uuid, type, status, trashed, title, area) VALUES ('P1', ?, ?, 0, 'Project One', 'A1');`, TaskTypeProject, StatusIncomplete); err != nil {
		return err
	}
	if _, err := conn.Exec(`INSERT INTO TMTask (uuid, type, status, trashed, title) VALUES ('P2', ?, ?, 0, 'Project Done');`, TaskTypeProject, StatusCompleted); err != nil {
		return err
	}
	if _, err := conn.Exec(`INSERT INTO TMTask (uuid, type, status, trashed, title, project, area, heading, notes) VALUES ('H1', ?, ?, 0, 'Heading', 'P1', 'A1', '', '');`, TaskTypeHeading, StatusIncomplete); err != nil {
		return err
	}
	if _, err := conn.Exec(`INSERT INTO TMTask (uuid, type, status, trashed, title, project, area, heading, notes, start, startDate, deadline, creationDate, userModificationDate, todayIndexReferenceDate) VALUES ('T1', ?, ?, 0, 'Task One', 'P1', 'A1', 'H1', 'Some notes', 1, ?, ?, ?, ?, 135004288);`, TaskTypeTodo, StatusIncomplete, startDate, deadline, nowUnix, nowUnix); err != nil {
		return err
	}
	if _, err := conn.Exec(`INSERT INTO TMTask (uuid, type, status, trashed, title, project, area, heading, start, startDate, deadline, creationDate, userModificationDate, rt1_recurrenceRule) VALUES ('T2', ?, ?, 0, 'Repeat Task', 'P1', 'A1', 'H1', 1, ?, ?, ?, ?, X'01');`, TaskTypeTodo, StatusIncomplete, startDate, deadline, nowUnix, nowUnix); err != nil {
		return err
	}
	if _, err := conn.Exec(`INSERT INTO TMTask (uuid, type, status, trashed, title, project, area, heading, start, startDate, deadline, creationDate, userModificationDate, rt1_repeatingTemplate) VALUES ('T3', ?, ?, 0, 'Generated Repeat Instance', 'P1', 'A1', 'H1', 1, ?, ?, ?, ?, 'T2');`, TaskTypeTodo, StatusCompleted, startDate, deadline, nowUnix, nowUnix); err != nil {
		return err
	}
	if _, err := conn.Exec(`INSERT INTO TMTag (uuid, title) VALUES ('TAG1', 'urgent');`); err != nil {
		return err
	}
	if _, err := conn.Exec(`INSERT INTO TMTaskTag (tasks, tags) VALUES ('T1', 'TAG1');`); err != nil {
		return err
	}
	if _, err := conn.Exec(`INSERT INTO TMChecklistItem (uuid, title, status, "index", task) VALUES ('C1', 'Check Item', ?, 0, 'T1');`, StatusIncomplete); err != nil {
		return err
	}
	return nil
}

func thingsDateForTest(t time.Time) int {
	date := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return date.Year()<<16 | int(date.Month())<<12 | date.Day()<<7
}

func TestNormalizeForMatch(t *testing.T) {
	cases := map[string]string{
		"💼 Rebtech":         "rebtech",
		"Rebtech":           "rebtech",
		"  Rebtech  ":       "rebtech",
		"🤳🍽️ PlateSnap":     "platesnap",
		"Rebtech Hemsida ":  "rebtech hemsida",
		"💰Personal Finance": "personal finance",
		"💼":                 "",
		"":                  "",
	}
	for in, want := range cases {
		if got := normalizeForMatch(in); got != want {
			t.Errorf("normalizeForMatch(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestResolveAreaIDEmojiInsensitive(t *testing.T) {
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer conn.Close()

	if _, err := conn.Exec(`CREATE TABLE TMArea (uuid TEXT PRIMARY KEY, title TEXT, visible INTEGER, "index" INTEGER);`); err != nil {
		t.Fatalf("create table: %v", err)
	}
	seed := []struct{ uuid, title string }{
		{"AR1", "💼 Rebtech"},
		{"AR2", "🤳🍽️ PlateSnap"},
		{"AR3", "🧪 Rebtech Lab"},
	}
	for _, a := range seed {
		if _, err := conn.Exec(`INSERT INTO TMArea (uuid, title, visible, "index") VALUES (?, ?, 1, 0);`, a.uuid, a.title); err != nil {
			t.Fatalf("insert area: %v", err)
		}
	}

	t.Run("exact title still resolves", func(t *testing.T) {
		got, err := resolveAreaID(conn, "💼 Rebtech")
		if err != nil || got != "AR1" {
			t.Fatalf("got (%q, %v), want AR1", got, err)
		}
	})

	t.Run("uuid still resolves", func(t *testing.T) {
		got, err := resolveAreaID(conn, "AR2")
		if err != nil || got != "AR2" {
			t.Fatalf("got (%q, %v), want AR2", got, err)
		}
	})

	t.Run("emoji-stripped exact resolves", func(t *testing.T) {
		got, err := resolveAreaID(conn, "PlateSnap")
		if err != nil || got != "AR2" {
			t.Fatalf("got (%q, %v), want AR2", got, err)
		}
	})

	t.Run("case-insensitive substring resolves", func(t *testing.T) {
		// "rebtech" is a substring of both "Rebtech" and "Rebtech Lab", but the
		// normalized-exact tier prefers the unique exact hit (AR1).
		got, err := resolveAreaID(conn, "rebtech")
		if err != nil || got != "AR1" {
			t.Fatalf("got (%q, %v), want AR1", got, err)
		}
	})

	t.Run("ambiguous substring errors with candidates", func(t *testing.T) {
		_, err := resolveAreaID(conn, "reb")
		if err == nil {
			t.Fatal("expected ambiguity error, got nil")
		}
		msg := err.Error()
		if !strings.Contains(msg, "ambiguous") || !strings.Contains(msg, "Rebtech") || !strings.Contains(msg, "Rebtech Lab") {
			t.Fatalf("unexpected error: %q", msg)
		}
	})

	t.Run("no match errors not found", func(t *testing.T) {
		_, err := resolveAreaID(conn, "Nonexistent")
		if err == nil || !strings.Contains(err.Error(), "not found") {
			t.Fatalf("expected not-found error, got %v", err)
		}
	})
}

func TestResolveProjectIDEmojiInsensitive(t *testing.T) {
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer conn.Close()

	if _, err := conn.Exec(`CREATE TABLE TMTask (uuid TEXT PRIMARY KEY, type INTEGER, status INTEGER, trashed INTEGER, title TEXT);`); err != nil {
		t.Fatalf("create table: %v", err)
	}
	rows := []struct {
		uuid, title       string
		taskType, trashed int
	}{
		{"PR1", "🚀 Launch Plan", TaskTypeProject, 0},
		{"PR2", "💼 Rebtech Hemsida", TaskTypeProject, 0},
		{"PR3", "🗑️ Launch Plan", TaskTypeProject, 1},      // trashed: must not fuzzy-match
		{"HD1", "Launch Plan Heading", TaskTypeHeading, 0}, // wrong type: must not match
		{"TD1", "Launch the rocket", TaskTypeTodo, 0},      // wrong type: must not match
	}
	for _, r := range rows {
		if _, err := conn.Exec(`INSERT INTO TMTask (uuid, type, status, trashed, title) VALUES (?, ?, 0, ?, ?);`, r.uuid, r.taskType, r.trashed, r.title); err != nil {
			t.Fatalf("insert task: %v", err)
		}
	}

	t.Run("emoji-stripped exact resolves the project only", func(t *testing.T) {
		// "Launch Plan" exists as an active project (PR1), a trashed project
		// (PR3, excluded), and a heading (HD1, wrong type) — only PR1 should win.
		got, err := resolveProjectID(conn, "Launch Plan")
		if err != nil || got != "PR1" {
			t.Fatalf("got (%q, %v), want PR1", got, err)
		}
	})

	t.Run("substring resolves across emoji", func(t *testing.T) {
		got, err := resolveProjectID(conn, "rebtech")
		if err != nil || got != "PR2" {
			t.Fatalf("got (%q, %v), want PR2", got, err)
		}
	})

	t.Run("trashed project is not a fuzzy match", func(t *testing.T) {
		// Drop the active project so only the trashed "Launch Plan" remains;
		// it must not resolve.
		if _, err := conn.Exec(`DELETE FROM TMTask WHERE uuid = 'PR1';`); err != nil {
			t.Fatalf("delete: %v", err)
		}
		_, err := resolveProjectID(conn, "Launch Plan")
		if err == nil || !strings.Contains(err.Error(), "not found") {
			t.Fatalf("expected not-found error, got %v", err)
		}
	})
}

func TestResolveTagIDEmojiInsensitive(t *testing.T) {
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer conn.Close()

	if _, err := conn.Exec(`CREATE TABLE TMTag (uuid TEXT PRIMARY KEY, title TEXT, shortcut TEXT, parent TEXT);`); err != nil {
		t.Fatalf("create table: %v", err)
	}
	tags := []struct{ uuid, title string }{
		{"TG1", "🔥 Urgent"},
		{"TG2", "⏳ Waiting"},
	}
	for _, tg := range tags {
		if _, err := conn.Exec(`INSERT INTO TMTag (uuid, title) VALUES (?, ?);`, tg.uuid, tg.title); err != nil {
			t.Fatalf("insert tag: %v", err)
		}
	}

	t.Run("emoji-stripped exact resolves", func(t *testing.T) {
		got, err := resolveTagID(conn, "urgent")
		if err != nil || got != "TG1" {
			t.Fatalf("got (%q, %v), want TG1", got, err)
		}
	})

	t.Run("no match errors not found", func(t *testing.T) {
		_, err := resolveTagID(conn, "nope")
		if err == nil || !strings.Contains(err.Error(), "not found") {
			t.Fatalf("expected not-found error, got %v", err)
		}
	})
}
