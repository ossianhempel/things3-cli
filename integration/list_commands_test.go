package integration_test

import (
	"database/sql"
	"strings"
	"testing"
	"time"
)

func TestInboxCommand(t *testing.T) {
	dbPath := writeTestDB(t)
	out, _, code := runThings(t, "", "inbox", "--db", dbPath)
	requireSuccess(t, code)
	assertContains(t, out, "Inbox Task")
}

func TestAnytimeCommand(t *testing.T) {
	dbPath := writeTestDB(t)
	out, _, code := runThings(t, "", "anytime", "--db", dbPath)
	requireSuccess(t, code)
	assertContains(t, out, "Anytime Task")
}

func TestSomedayCommand(t *testing.T) {
	dbPath := writeTestDB(t)
	out, _, code := runThings(t, "", "someday", "--db", dbPath)
	requireSuccess(t, code)
	assertContains(t, out, "Someday Task")
}

func TestUpcomingCommand(t *testing.T) {
	dbPath := writeTestDB(t)
	out, _, code := runThings(t, "", "upcoming", "--db", dbPath)
	requireSuccess(t, code)
	assertContains(t, out, "Upcoming Task")
}

func TestTemplatesCommand(t *testing.T) {
	dbPath := writeTestDB(t)
	out, _, code := runThings(t, "", "templates", "--db", dbPath, "--area", "Home")
	requireSuccess(t, code)
	assertContains(t, out, "Template Task")
	assertNotContains(t, out, "Generated Template Instance")
}

func TestDeadlinesCommand(t *testing.T) {
	dbPath := writeTestDB(t)
	out, _, code := runThings(t, "", "deadlines", "--db", dbPath)
	requireSuccess(t, code)
	assertContains(t, out, "Deadline Task")
}

func TestCompletedCommand(t *testing.T) {
	dbPath := writeTestDB(t)
	out, _, code := runThings(t, "", "completed", "--db", dbPath)
	requireSuccess(t, code)
	assertContains(t, out, "Completed Task")
}

func TestCanceledCommand(t *testing.T) {
	dbPath := writeTestDB(t)
	out, _, code := runThings(t, "", "canceled", "--db", dbPath)
	requireSuccess(t, code)
	assertContains(t, out, "Canceled Task")
}

func TestTrashCommand(t *testing.T) {
	dbPath := writeTestDB(t)
	out, _, code := runThings(t, "", "trash", "--db", dbPath)
	requireSuccess(t, code)
	assertContains(t, out, "Trashed Task")
}

func TestLogbookCommand(t *testing.T) {
	dbPath := writeTestDB(t)
	out, _, code := runThings(t, "", "logbook", "--db", dbPath)
	requireSuccess(t, code)
	assertContains(t, out, "Completed Task")
}

func TestLogTodayCommand(t *testing.T) {
	dbPath := writeTestDB(t)
	out, _, code := runThings(t, "", "logtoday", "--db", dbPath)
	requireSuccess(t, code)
	assertContains(t, out, "Completed Task")
}

func TestCreatedTodayCommand(t *testing.T) {
	dbPath := writeTestDB(t)
	out, _, code := runThings(t, "", "createdtoday", "--db", dbPath)
	requireSuccess(t, code)
	assertContains(t, out, "Task One")
}

func TestAllCommand(t *testing.T) {
	dbPath := writeTestDB(t)
	out, _, code := runThings(t, "", "all", "--db", dbPath)
	requireSuccess(t, code)
	assertContains(t, out, "Inbox")
	assertContains(t, out, "Inbox Task")
}

func TestTodayCompositeOrderAndExplicitSort(t *testing.T) {
	dbPath := writeTestDB(t)
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	today := thingsDate(time.Now())
	for _, row := range []struct {
		uuid, title              string
		bucket, reference, index int
	}{
		{"ORDER_NEW", "Zulu New Reference", 0, 300, 9},
		{"ORDER_OLD", "Alpha Old Reference", 0, 200, 1},
		{"ORDER_EVENING", "Beta Evening", 1, 999, 0},
	} {
		if _, err := conn.Exec(`INSERT INTO TMTask (uuid, type, status, trashed, title, start, startDate, startBucket, todayIndexReferenceDate, todayIndex) VALUES (?, 0, 0, 0, ?, 1, ?, ?, ?, ?)`, row.uuid, row.title, today, row.bucket, row.reference, row.index); err != nil {
			t.Fatalf("insert %s: %v", row.uuid, err)
		}
	}
	if err := conn.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}

	out, _, code := runThings(t, "", "today", "--db", dbPath, "--format", "csv", "--select", "title", "--no-header")
	requireSuccess(t, code)
	assertBefore(t, out, "Zulu New Reference", "Alpha Old Reference")
	assertBefore(t, out, "Alpha Old Reference", "Beta Evening")

	out, _, code = runThings(t, "", "today", "--db", dbPath, "--format", "csv", "--select", "title", "--no-header", "--sort", "title")
	requireSuccess(t, code)
	assertBefore(t, out, "Alpha Old Reference", "Zulu New Reference")

	out, _, code = runThings(t, "", "all", "--db", dbPath, "--json")
	requireSuccess(t, code)
	assertBefore(t, out, "Zulu New Reference", "Alpha Old Reference")
	assertBefore(t, out, "Alpha Old Reference", "Beta Evening")
}

func assertBefore(t *testing.T, output, first, second string) {
	t.Helper()
	firstIndex := strings.Index(output, first)
	secondIndex := strings.Index(output, second)
	if firstIndex < 0 || secondIndex < 0 || firstIndex >= secondIndex {
		t.Fatalf("expected %q before %q in %q", first, second, output)
	}
}

func TestAreasRecursive(t *testing.T) {
	dbPath := writeTestDB(t)
	out, _, code := runThings(t, "", "areas", "--db", dbPath, "--recursive")
	requireSuccess(t, code)
	assertContains(t, out, "Project One")
	assertContains(t, out, "Task One")
}

func TestAreasOnlyProjects(t *testing.T) {
	dbPath := writeTestDB(t)
	out, _, code := runThings(t, "", "areas", "--db", dbPath, "--recursive", "--only-projects")
	requireSuccess(t, code)
	assertContains(t, out, "Project One")
	assertNotContains(t, out, "Task One")
}

func TestProjectsRecursive(t *testing.T) {
	dbPath := writeTestDB(t)
	out, _, code := runThings(t, "", "projects", "--db", dbPath, "--recursive")
	requireSuccess(t, code)
	assertContains(t, out, "Project One")
	assertContains(t, out, "Task One")
}

func TestProjectsOnlyProjects(t *testing.T) {
	dbPath := writeTestDB(t)
	out, _, code := runThings(t, "", "projects", "--db", dbPath, "--recursive", "--only-projects")
	requireSuccess(t, code)
	assertContains(t, out, "Project One")
	assertNotContains(t, out, "Task One")
}

func TestTodayCommandExcludesOpenTasksInClosedProjects(t *testing.T) {
	dbPath := writeTestDB(t)
	out, _, code := runThings(t, "", "today", "--db", dbPath)
	requireSuccess(t, code)
	assertContains(t, out, "Today Task")
	assertNotContains(t, out, "Open Task in Completed Project")
	assertNotContains(t, out, "Scheduled Task in Completed Project")
}

func TestTasksCommandExcludesOpenTasksInClosedProjects(t *testing.T) {
	dbPath := writeTestDB(t)
	out, _, code := runThings(t, "", "tasks", "--db", dbPath, "--limit", "0")
	requireSuccess(t, code)
	assertContains(t, out, "Task One")
	assertNotContains(t, out, "Open Task in Completed Project")
	assertNotContains(t, out, "Open Task in Canceled Project")
}

func TestSearchExcludesOpenTasksInClosedProjects(t *testing.T) {
	dbPath := writeTestDB(t)
	out, _, code := runThings(t, "", "search", "Completed Project", "--db", dbPath)
	requireSuccess(t, code)
	assertNotContains(t, out, "Open Task in Completed Project")
}

func TestCompletedCommandKeepsClosedProjectChildren(t *testing.T) {
	dbPath := writeTestDB(t)
	out, _, code := runThings(t, "", "completed", "--db", dbPath)
	requireSuccess(t, code)
	assertContains(t, out, "Completed Task in Completed Project")
}

func TestShowResolvesOpenTaskInClosedProject(t *testing.T) {
	dbPath := writeTestDB(t)
	out, _, code := runThings(t, "", "show", "--db", dbPath, "--id", "LEFTTODAY1")
	requireSuccess(t, code)
	assertContains(t, out, "Open Task in Completed Project")
}
