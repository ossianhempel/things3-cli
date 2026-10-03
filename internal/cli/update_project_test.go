package cli

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestUpdateProjectCommandRequiresAuthToken(t *testing.T) {
	t.Setenv("THINGS_AUTH_TOKEN", "")
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"update-project", "--id", "123", "Title"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err == nil {
		t.Fatalf("expected error")
	}
	if len(launcher.args) != 0 {
		t.Fatalf("expected no open invocation")
	}
}

func TestUpdateProjectCommandWithAuthAndID(t *testing.T) {
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"update-project", "--auth-token", "tok", "--id", "123", "Title"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	url := requireOpenURL(t, launcher)
	if !strings.Contains(url, "auth-token=tok") {
		t.Fatalf("expected auth-token in url, got %q", url)
	}
	if !strings.Contains(url, "id=123") {
		t.Fatalf("expected id in url, got %q", url)
	}
}

func TestUpdateProjectCommandEmptyValuesClearFields(t *testing.T) {
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"update-project", "--auth-token", "tok", "--id", "123", "--when=", "--deadline=", "--tags="})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	url := requireOpenURL(t, launcher)
	for _, param := range []string{"when=&", "deadline=&", "tags=&"} {
		if !strings.Contains(url, param) {
			t.Fatalf("expected %q in url, got %q", param, url)
		}
	}
}

func TestUpdateProjectCommandAddTags(t *testing.T) {
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"update-project", "--auth-token", "tok", "--id", "123", "--add-tags", "Focus,Home"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	url := requireOpenURL(t, launcher)
	if !strings.Contains(url, "add-tags=Focus%2CHome") {
		t.Fatalf("expected add-tags in url, got %q", url)
	}
}

func TestRepeatUpdateProjectJSONDryRunResolvesTargetWithoutLaunch(t *testing.T) {
	dbPath := writeTestDB(t)
	launcher := &recordLauncher{}
	out := &bytes.Buffer{}
	app := &App{In: strings.NewReader(""), Out: out, Err: &bytes.Buffer{}, Launcher: launcher, DryRun: true}
	root := NewRoot(app)
	root.SetArgs([]string{"--dry-run", "update-project", "--db", dbPath, "--id", "P1", "--repeat=week", "--repeat-start=2025-01-02", "--json"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if len(launcher.args) != 0 {
		t.Fatalf("dry-run launched Things: %#v", launcher.args)
	}
	var result repeatResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v: %q", err, out.String())
	}
	if result.IDs.Requested != "P1" || result.IDs.Template != "P1" || len(result.Stages) != 2 {
		t.Fatalf("unexpected preview: %#v", result)
	}
}

func TestRepeatUpdateProjectAppliesAndVerifiesTemplate(t *testing.T) {
	dbPath := writeTestDB(t)
	out := &bytes.Buffer{}
	app := &App{In: strings.NewReader(""), Out: out, Err: &bytes.Buffer{}, Launcher: &recordLauncher{}}
	root := NewRoot(app)
	root.SetArgs([]string{"update-project", "--db", dbPath, "--id", "P1", "--repeat=week", "--repeat-mode=schedule", "--repeat-start=2025-01-02", "--json"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	var result repeatResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v: %q", err, out.String())
	}
	if !result.Verified || result.IDs.Template != "P1" || !result.Repeat.Active {
		t.Fatalf("unexpected result: %#v", result)
	}
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	var rule []byte
	var start int
	if err := conn.QueryRow(`SELECT rt1_recurrenceRule, start FROM TMTask WHERE uuid = 'P1'`).Scan(&rule, &start); err != nil {
		t.Fatal(err)
	}
	if len(rule) == 0 || start != 2 {
		t.Fatalf("project not made repeating: rule=%d start=%d", len(rule), start)
	}
}

func TestRepeatUpdateProjectClearRemovesRule(t *testing.T) {
	dbPath := writeTestDB(t)
	out := &bytes.Buffer{}
	app := &App{In: strings.NewReader(""), Out: out, Err: &bytes.Buffer{}, Launcher: &recordLauncher{}}
	root := NewRoot(app)
	root.SetArgs([]string{"update-project", "--db", dbPath, "--id", "P1", "--repeat=week", "--repeat-mode=schedule", "--repeat-start=2025-01-02", "--json"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	root = NewRoot(app)
	root.SetArgs([]string{"update-project", "--db", dbPath, "--id", "P1", "--repeat-clear", "--json"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	var result repeatResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v: %q", err, out.String())
	}
	if !result.Verified || result.Action != "clear" {
		t.Fatalf("unexpected result: %#v", result)
	}
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	var rule []byte
	if err := conn.QueryRow(`SELECT rt1_recurrenceRule FROM TMTask WHERE uuid = 'P1'`).Scan(&rule); err != nil {
		t.Fatal(err)
	}
	if len(rule) != 0 {
		t.Fatalf("rule not cleared: %d bytes", len(rule))
	}
}

func TestRepeatUpdateProjectRejectsTodoTarget(t *testing.T) {
	dbPath := writeTestDB(t)
	app := &App{In: strings.NewReader(""), Out: &bytes.Buffer{}, Err: &bytes.Buffer{}, Launcher: &recordLauncher{}}
	root := NewRoot(app)
	root.SetArgs([]string{"update-project", "--db", dbPath, "--id", "T1", "--repeat=day"})
	if err := root.Execute(); err == nil || !strings.Contains(err.Error(), "item type mismatch") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRepeatUpdateProjectRejectsStatusChanges(t *testing.T) {
	for _, arg := range []string{"--completed", "--canceled"} {
		app := &App{In: strings.NewReader(""), Out: &bytes.Buffer{}, Err: &bytes.Buffer{}, Launcher: &recordLauncher{}}
		root := NewRoot(app)
		root.SetArgs([]string{"update-project", "--id", "P1", "--repeat=day", arg})
		if err := root.Execute(); err == nil || !strings.Contains(err.Error(), "cannot be combined") {
			t.Fatalf("unexpected error for %s: %v", arg, err)
		}
	}
}

func TestRepeatUpdateProjectRejectsDuplicate(t *testing.T) {
	launcher := &recordLauncher{}
	app := &App{In: strings.NewReader(""), Out: &bytes.Buffer{}, Err: &bytes.Buffer{}, Launcher: launcher}
	root := NewRoot(app)
	root.SetArgs([]string{"update-project", "--id", "P1", "--repeat=day", "--duplicate"})
	if err := root.Execute(); err == nil || !strings.Contains(err.Error(), "--duplicate cannot be combined") {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(launcher.args) != 0 {
		t.Fatalf("duplicate+repeat launched Things: %#v", launcher.args)
	}
}

func TestRepeatUpdateProjectRejectsCountNotGreaterThanExisting(t *testing.T) {
	dbPath := writeTestDB(t)
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(`UPDATE TMTask SET rt1_instanceCreationCount = 5, rt1_recurrenceRule = X'01', start = 2, startBucket = 0 WHERE uuid = 'P1'`); err != nil {
		t.Fatal(err)
	}
	conn.Close()
	for _, count := range []int{3, 5} {
		app := &App{In: strings.NewReader(""), Out: &bytes.Buffer{}, Err: &bytes.Buffer{}, Launcher: &recordLauncher{}}
		root := NewRoot(app)
		root.SetArgs([]string{"update-project", "--db", dbPath, "--id", "P1", "--repeat=day", "--repeat-count", fmt.Sprint(count)})
		if err := root.Execute(); err == nil || !strings.Contains(err.Error(), "not greater than the 5 occurrences") {
			t.Fatalf("count %d: unexpected error: %v", count, err)
		}
	}
}

func TestRepeatUpdateProjectAcceptsCountGreaterThanExisting(t *testing.T) {
	dbPath := writeTestDB(t)
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(`UPDATE TMTask SET rt1_instanceCreationCount = 5, rt1_recurrenceRule = X'01', start = 2, startBucket = 0 WHERE uuid = 'P1'`); err != nil {
		t.Fatal(err)
	}
	conn.Close()
	out := &bytes.Buffer{}
	app := &App{In: strings.NewReader(""), Out: out, Err: &bytes.Buffer{}, Launcher: &recordLauncher{}}
	root := NewRoot(app)
	root.SetArgs([]string{"update-project", "--db", dbPath, "--id", "P1", "--repeat=day", "--repeat-count=20", "--json"})
	if err := root.Execute(); err != nil {
		t.Fatalf("total greater than existing count should apply: %v", err)
	}
	var result repeatResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if !result.Verified || result.Repeat.Count == nil || *result.Repeat.Count != 20 {
		t.Fatalf("expected verified count 20: %#v", result.Repeat)
	}
}

func TestRepeatUpdateProjectDatabaseFailureEmitsResult(t *testing.T) {
	dbPath := writeTestDB(t)
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	_, err = conn.Exec(`CREATE TRIGGER fail_repeat BEFORE UPDATE OF rt1_recurrenceRule ON TMTask BEGIN SELECT RAISE(FAIL, 'forced repeat failure'); END;`)
	conn.Close()
	if err != nil {
		t.Fatal(err)
	}
	out := &bytes.Buffer{}
	app := &App{In: strings.NewReader(""), Out: out, Err: &bytes.Buffer{}, Launcher: &recordLauncher{}}
	root := NewRoot(app)
	root.SetArgs([]string{"update-project", "--db", dbPath, "--id", "P1", "--repeat=week", "--repeat-start=2025-01-02", "--json"})
	if err := root.Execute(); err == nil {
		t.Fatal("expected failure")
	}
	var result repeatResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("missing result: %v %q", err, out.String())
	}
	if len(result.Stages) != 1 || result.Stages[0].Status != "failed" || len(result.Recovery) != 1 {
		t.Fatalf("unexpected result: %#v", result)
	}
	if result.Recovery[0].Argv[1] != "update-project" {
		t.Fatalf("recovery should target update-project: %#v", result.Recovery[0].Argv)
	}
}

