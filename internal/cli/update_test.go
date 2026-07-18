package cli

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"strings"
	"testing"
)

type triggerLauncher struct{ dbPath string }

func (l *triggerLauncher) Open(args ...string) error {
	conn, err := sql.Open("sqlite", l.dbPath)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Exec(`CREATE TRIGGER fail_repeat BEFORE UPDATE OF rt1_recurrenceRule ON TMTask BEGIN SELECT RAISE(FAIL, 'forced repeat failure'); END;`)
	return err
}

func TestRepeatUpdateJSONDryRunResolvesTargetWithoutLaunch(t *testing.T) {
	dbPath := writeTestDB(t)
	launcher := &recordLauncher{}
	out := &bytes.Buffer{}
	app := &App{In: strings.NewReader(""), Out: out, Err: &bytes.Buffer{}, Launcher: launcher, DryRun: true}
	root := NewRoot(app)
	root.SetArgs([]string{"--dry-run", "update", "--db", dbPath, "--id", "T1", "--repeat=week", "--repeat-start=2025-01-02", "--json"})
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
	if result.IDs.Requested != "T1" || result.IDs.Template != "T1" || len(result.Stages) != 2 {
		t.Fatalf("unexpected preview: %#v", result)
	}
}

func TestRepeatUpdateDryRunRedactsURLIntent(t *testing.T) {
	dbPath := writeTestDB(t)
	out := &bytes.Buffer{}
	app := &App{In: strings.NewReader(""), Out: out, Err: &bytes.Buffer{}, Launcher: &recordLauncher{}, DryRun: true}
	root := NewRoot(app)
	root.SetArgs([]string{"--dry-run", "update", "--db", dbPath, "--id", "T1", "--auth-token", "secret", "--notes", "changed", "--repeat=day", "--json"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	var result repeatResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Intent == nil || strings.Contains(result.Intent.URL, "secret") || !strings.Contains(result.Intent.URL, "REDACTED") {
		t.Fatalf("unsafe intent: %#v", result.Intent)
	}
}

func TestRepeatOnlyDatabaseFailureStillEmitsResult(t *testing.T) {
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
	root.SetArgs([]string{"update", "--db", dbPath, "--id", "T1", "--repeat=week", "--repeat-start=2025-01-02", "--json"})
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
}

func TestCombinedUpdateDatabaseFailureIsPartial(t *testing.T) {
	dbPath := writeTestDB(t)
	out := &bytes.Buffer{}
	app := &App{In: strings.NewReader(""), Out: out, Err: &bytes.Buffer{}, Launcher: &triggerLauncher{dbPath: dbPath}}
	root := NewRoot(app)
	root.SetArgs([]string{"update", "--db", dbPath, "--id", "T1", "--auth-token", "tok", "--notes", "changed", "--repeat=day", "--repeat-start=2025-01-02", "--json"})
	if err := root.Execute(); err == nil {
		t.Fatal("expected failure")
	}
	var result repeatResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if !result.Partial || result.Stages[0].Name != "url" || result.Stages[1].Status != "failed" {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestRepeatUpdateRejectsWhenAndLater(t *testing.T) {
	for _, arg := range []string{"--when=today", "--later"} {
		app := &App{In: strings.NewReader(""), Out: &bytes.Buffer{}, Err: &bytes.Buffer{}, Launcher: &recordLauncher{}}
		root := NewRoot(app)
		root.SetArgs([]string{"update", "--id", "T1", "--repeat=day", arg})
		if err := root.Execute(); err == nil || !strings.Contains(err.Error(), "cannot be combined") {
			t.Fatalf("%s: %v", arg, err)
		}
	}
}

func TestUpdateCommandRequiresAuthToken(t *testing.T) {
	t.Setenv("THINGS_AUTH_TOKEN", "")
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"update", "--id", "123", "Title"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err == nil {
		t.Fatalf("expected error")
	}
	if len(launcher.args) != 0 {
		t.Fatalf("expected no open invocation")
	}
}

func TestUpdateCommandWithAuthAndID(t *testing.T) {
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"update", "--auth-token", "tok", "--id", "123", "--no-verify", "Title"})
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

func TestUpdateCommandLaterFlag(t *testing.T) {
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"update", "--auth-token", "tok", "--id", "123", "--later", "--no-verify"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	url := requireOpenURL(t, launcher)
	if !strings.Contains(url, "when=evening") {
		t.Fatalf("expected when=evening in url, got %q", url)
	}
}

func TestUpdateCommandEmptyValuesClearFields(t *testing.T) {
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"update", "--auth-token", "tok", "--id", "123", "--when=", "--deadline=", "--tags=", "--no-verify"})
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

func TestUpdateCommandRejectsUnsafeTitle(t *testing.T) {
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"update", "--auth-token", "tok", "--id", "123", "--no-verify", "tag=work"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err == nil {
		t.Fatalf("expected error")
	}
	if len(launcher.args) != 0 {
		t.Fatalf("expected no open invocation")
	}
}

func TestUpdateCommandAllowsUnsafeTitleWithFlag(t *testing.T) {
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"update", "--auth-token", "tok", "--id", "123", "--allow-unsafe-title", "--no-verify", "tag=work"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	url := requireOpenURL(t, launcher)
	if !strings.Contains(url, "title=tag%3Dwork") {
		t.Fatalf("expected title in url, got %q", url)
	}
}

func TestUpdateCommandBlocksEveningForNonToday(t *testing.T) {
	dbPath := writeTestDB(t)
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"update", "--db", dbPath, "--auth-token", "tok", "--id", "UP1", "--when=evening", "--no-verify"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err == nil {
		t.Fatalf("expected error")
	}
	if len(launcher.args) != 0 {
		t.Fatalf("expected no open invocation")
	}
}

func TestUpdateCommandAllowsEveningForNonTodayWithFlag(t *testing.T) {
	dbPath := writeTestDB(t)
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"update", "--db", dbPath, "--auth-token", "tok", "--id", "UP1", "--when=evening", "--allow-non-today", "--no-verify"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	if len(launcher.args) == 0 {
		t.Fatalf("expected open invocation")
	}
}

func TestUpdateCommandCompletesChecklistItemViaJSON(t *testing.T) {
	dbPath := writeTestDB(t)
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"update", "--db", dbPath, "--auth-token", "tok", "--id", "T1", "--complete-checklist-item", "Check Item"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	url := requireOpenURL(t, launcher)
	if !strings.HasPrefix(url, "things:///json?") {
		t.Fatalf("expected json url, got %q", url)
	}
	if !strings.Contains(url, "auth-token=tok") {
		t.Fatalf("expected auth token in url, got %q", url)
	}
	if !strings.Contains(url, "completed%22%3Atrue") {
		t.Fatalf("expected completed checklist payload, got %q", url)
	}
}

func TestUpdateCommandChecklistStatusRequiresID(t *testing.T) {
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"update", "--auth-token", "tok", "--complete-checklist-item", "Check Item"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err == nil {
		t.Fatalf("expected error")
	}
	if len(launcher.args) != 0 {
		t.Fatalf("expected no open invocation")
	}
}

func TestUpdateCommandChecklistStatusRejectsOtherChanges(t *testing.T) {
	dbPath := writeTestDB(t)
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"update", "--db", dbPath, "--auth-token", "tok", "--id", "T1", "--complete-checklist-item", "Check Item", "--notes", "Nope"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err == nil {
		t.Fatalf("expected error")
	}
	if len(launcher.args) != 0 {
		t.Fatalf("expected no open invocation")
	}
}
