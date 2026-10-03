package db

import (
	"database/sql"
	"testing"
	"time"
)

// seedClosedContextDB builds a fixture where open to-dos are left behind inside
// projects the user has completed or canceled. Things treats those leftovers as
// archived and hides them from the active lists.
func seedClosedContextDB(t *testing.T) *Store {
	t.Helper()
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { conn.Close() })

	if err := seedTodayDB(conn); err != nil {
		t.Fatalf("seed db: %v", err)
	}
	if _, err := conn.Exec(`DELETE FROM TMTask`); err != nil {
		t.Fatalf("clear tasks: %v", err)
	}
	if _, err := conn.Exec(`CREATE TABLE TMChecklistItem (
		uuid TEXT PRIMARY KEY, userModificationDate REAL, creationDate REAL,
		title TEXT, status INTEGER, stopDate REAL, "index" INTEGER, task TEXT);`); err != nil {
		t.Fatalf("create checklist table: %v", err)
	}

	today := thingsDate(time.Now())
	yesterday := thingsDate(time.Now().AddDate(0, 0, -1))

	projects := []struct {
		uuid   string
		title  string
		status int
	}{
		{"P_OPEN", "Open Project", StatusIncomplete},
		{"P_DONE", "Completed Project", StatusCompleted},
		{"P_CANCELED", "Canceled Project", StatusCanceled},
	}
	for _, project := range projects {
		if _, err := conn.Exec(
			`INSERT INTO TMTask (uuid, type, status, trashed, title, start, startBucket)
			 VALUES (?, ?, ?, 0, ?, 1, 0)`,
			project.uuid, TaskTypeProject, project.status, project.title); err != nil {
			t.Fatalf("insert project %s: %v", project.uuid, err)
		}
	}

	todos := []struct {
		uuid      string
		title     string
		project   string
		status    int
		start     int
		startDate int
	}{
		// Legitimately in Today.
		{"OPEN_IN_OPEN", "Open todo in open project", "P_OPEN", StatusIncomplete, 1, today},
		// Left open when the parent project was closed. Things hides these.
		{"OPEN_IN_DONE", "Open todo in completed project", "P_DONE", StatusIncomplete, 1, today},
		{"OPEN_IN_CANCELED", "Open todo in canceled project", "P_CANCELED", StatusIncomplete, 1, today},
		// Scheduled in the past, which is how these usually reach Today.
		{"OPEN_SCHEDULED_IN_DONE", "Open scheduled todo in completed project", "P_DONE", StatusIncomplete, 2, yesterday},
		// Closed alongside the project. Must stay visible in the Logbook.
		{"DONE_IN_DONE", "Completed todo in completed project", "P_DONE", StatusCompleted, 1, today},
	}
	for _, todo := range todos {
		if _, err := conn.Exec(
			`INSERT INTO TMTask (uuid, type, status, trashed, title, project, start, startDate, startBucket, todayIndex)
			 VALUES (?, ?, ?, 0, ?, ?, ?, ?, 0, 0)`,
			todo.uuid, TaskTypeTodo, todo.status, todo.title, todo.project, todo.start, todo.startDate); err != nil {
			t.Fatalf("insert todo %s: %v", todo.uuid, err)
		}
	}

	return &Store{conn: conn, path: ":memory:"}
}

func closedContextUUIDs(tasks []Task) map[string]bool {
	seen := make(map[string]bool, len(tasks))
	for _, task := range tasks {
		seen[task.UUID] = true
	}
	return seen
}

func TestTodayExcludesOpenTasksInClosedProjects(t *testing.T) {
	store := seedClosedContextDB(t)
	status := StatusIncomplete

	tasks, err := store.TodayTasks(TaskFilter{Status: &status, ExcludeTrashedContext: true})
	if err != nil {
		t.Fatalf("today tasks: %v", err)
	}
	got := closedContextUUIDs(tasks)

	if !got["OPEN_IN_OPEN"] {
		t.Errorf("expected OPEN_IN_OPEN in Today, got %v", got)
	}
	for _, unwanted := range []string{"OPEN_IN_DONE", "OPEN_IN_CANCELED", "OPEN_SCHEDULED_IN_DONE"} {
		if got[unwanted] {
			t.Errorf("open todo %s belongs to a closed project and must not appear in Today", unwanted)
		}
	}
}

func TestTasksExcludeOpenTasksInClosedProjects(t *testing.T) {
	store := seedClosedContextDB(t)
	status := StatusIncomplete

	tasks, err := store.Tasks(TaskFilter{Status: &status, ExcludeTrashedContext: true})
	if err != nil {
		t.Fatalf("tasks: %v", err)
	}
	got := closedContextUUIDs(tasks)

	if !got["OPEN_IN_OPEN"] {
		t.Errorf("expected OPEN_IN_OPEN in tasks, got %v", got)
	}
	for _, unwanted := range []string{"OPEN_IN_DONE", "OPEN_IN_CANCELED", "OPEN_SCHEDULED_IN_DONE"} {
		if got[unwanted] {
			t.Errorf("open todo %s belongs to a closed project and must not appear in tasks", unwanted)
		}
	}
}

// Completing a project completes its children too, and the Logbook must keep
// showing them, so the closed-context guard must not reach those rows.
func TestCompletedTasksKeepClosedProjectChildren(t *testing.T) {
	store := seedClosedContextDB(t)

	tasks, err := store.CompletedTasks(TaskFilter{ExcludeTrashedContext: true})
	if err != nil {
		t.Fatalf("completed tasks: %v", err)
	}
	if got := closedContextUUIDs(tasks); !got["DONE_IN_DONE"] {
		t.Errorf("completed child of a completed project must stay in the Logbook, got %v", got)
	}
}

// An explicit project filter treats a closed parent exactly like the existing
// code already treats a trashed parent: the open leftovers stay hidden, and
// show --id is the escape hatch. This pins the two guards to the same
// behavior so they cannot drift apart.
func TestExplicitProjectFilterMatchesTrashedContextBehavior(t *testing.T) {
	store := seedClosedContextDB(t)
	if _, err := store.conn.Exec(
		`INSERT INTO TMTask (uuid, type, status, trashed, title, start, startBucket)
		 VALUES ('P_TRASHED', ?, ?, 1, 'Trashed Project', 1, 0)`,
		TaskTypeProject, StatusIncomplete); err != nil {
		t.Fatalf("insert trashed project: %v", err)
	}
	if _, err := store.conn.Exec(
		`INSERT INTO TMTask (uuid, type, status, trashed, title, project, start, startBucket, todayIndex)
		 VALUES ('OPEN_IN_TRASHED', ?, ?, 0, 'Open todo in trashed project', 'P_TRASHED', 1, 0, 0)`,
		TaskTypeTodo, StatusIncomplete); err != nil {
		t.Fatalf("insert trashed project child: %v", err)
	}

	status := StatusIncomplete
	for _, project := range []string{"P_TRASHED", "P_DONE", "P_CANCELED"} {
		tasks, err := store.Tasks(TaskFilter{
			Status:                &status,
			ExcludeTrashedContext: true,
			ProjectID:             project,
		})
		if err != nil {
			t.Fatalf("tasks for %s: %v", project, err)
		}
		if len(tasks) != 0 {
			t.Errorf("explicit filter on closed project %s returned %d open rows, want 0", project, len(tasks))
		}
	}
}

// show --id must still resolve a leftover open to-do, otherwise the UUID-first
// agent workflow breaks for exactly these rows.
func TestTaskByIDStillResolvesClosedProjectChild(t *testing.T) {
	store := seedClosedContextDB(t)

	task, err := store.TaskByID("OPEN_IN_DONE")
	if err != nil {
		t.Fatalf("task by id: %v", err)
	}
	if task.UUID != "OPEN_IN_DONE" {
		t.Errorf("expected OPEN_IN_DONE, got %s", task.UUID)
	}
}
