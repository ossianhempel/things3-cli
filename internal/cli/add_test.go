package cli

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"howett.net/plist"
)

type insertingLauncher struct {
	dbPath string
	calls  int
}

func (l *insertingLauncher) Open(args ...string) error {
	l.calls++
	if len(args) == 0 || !strings.HasPrefix(args[len(args)-1], "things:///") {
		return nil
	}
	conn, err := sql.Open("sqlite", l.dbPath)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Exec(`INSERT INTO TMTask (uuid,type,status,trashed,title,start,startBucket,creationDate,userModificationDate) VALUES ('CREATED',0,0,0,'Repeat me',1,4,?,?)`, float64(time.Now().Unix()), float64(time.Now().Unix()))
	return err
}

func TestRepeatAddJSONDryRunDescribesAllStagesWithoutLaunch(t *testing.T) {
	dbPath := writeTestDB(t)
	launcher := &recordLauncher{}
	out := &bytes.Buffer{}
	app := &App{In: strings.NewReader(""), Out: out, Err: &bytes.Buffer{}, Launcher: launcher, DryRun: true}
	root := NewRoot(app)
	root.SetArgs([]string{"--dry-run", "add", "--db", dbPath, "--repeat=day", "--repeat-mode=schedule", "--repeat-start=2025-01-02", "--json", "Repeat me"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if len(launcher.args) != 0 {
		t.Fatalf("dry-run launched Things: %#v", launcher.args)
	}
	var result repeatResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("stdout is not one JSON object: %v: %q", err, out.String())
	}
	if !result.DryRun || result.Repeat.Mode != "schedule" || len(result.Stages) != 4 {
		t.Fatalf("unexpected preview: %#v", result)
	}
	if result.Database.Path == "" || result.Database.Source != "flag" || result.Intent == nil || !strings.HasPrefix(result.Intent.URL, "things:///add?") {
		t.Fatalf("missing preflight provenance/intent: %#v", result)
	}
}

func TestRepeatAddRejectsWhenBeforeLaunch(t *testing.T) {
	launcher := &recordLauncher{}
	app := &App{In: strings.NewReader(""), Out: &bytes.Buffer{}, Err: &bytes.Buffer{}, Launcher: launcher}
	root := NewRoot(app)
	root.SetArgs([]string{"add", "--repeat=day", "--when=today", "Nope"})
	if err := root.Execute(); err == nil || !strings.Contains(err.Error(), "cannot be combined") {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(launcher.args) != 0 {
		t.Fatal("launched Things")
	}
}

func TestRepeatAddInvalidSpecFailsBeforeLaunch(t *testing.T) {
	launcher := &recordLauncher{}
	app := &App{In: strings.NewReader(""), Out: &bytes.Buffer{}, Err: &bytes.Buffer{}, Launcher: launcher}
	root := NewRoot(app)
	root.SetArgs([]string{"add", "--dry-run", "--repeat=day", "--repeat-every=0", "Nope"})
	if err := root.Execute(); err == nil {
		t.Fatal("expected invalid interval")
	}
	if len(launcher.args) != 0 {
		t.Fatalf("invalid spec launched Things: %#v", launcher.args)
	}
}

func TestRepeatAddCountAndUntilAreMutuallyExclusive(t *testing.T) {
	launcher := &recordLauncher{}
	app := &App{In: strings.NewReader(""), Out: &bytes.Buffer{}, Err: &bytes.Buffer{}, Launcher: launcher}
	root := NewRoot(app)
	root.SetArgs([]string{"add", "--repeat=day", "--repeat-count=20", "--repeat-until=2026-12-24", "Nope"})
	if err := root.Execute(); err == nil || !strings.Contains(err.Error(), "mutually exclusive") {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(launcher.args) != 0 {
		t.Fatalf("invalid combo launched Things: %#v", launcher.args)
	}
}

func TestRepeatAddCountMustBePositive(t *testing.T) {
	launcher := &recordLauncher{}
	app := &App{In: strings.NewReader(""), Out: &bytes.Buffer{}, Err: &bytes.Buffer{}, Launcher: launcher}
	root := NewRoot(app)
	root.SetArgs([]string{"add", "--repeat=day", "--repeat-count=0", "Nope"})
	if err := root.Execute(); err == nil || !strings.Contains(err.Error(), "positive integer") {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(launcher.args) != 0 {
		t.Fatalf("invalid count launched Things: %#v", launcher.args)
	}
}

func TestRepeatAddAppliesAndVerifiesCreatedTemplate(t *testing.T) {
	dbPath := writeTestDB(t)
	launcher := &insertingLauncher{dbPath: dbPath}
	out := &bytes.Buffer{}
	app := &App{In: strings.NewReader(""), Out: out, Err: &bytes.Buffer{}, Launcher: launcher}
	root := NewRoot(app)
	root.SetArgs([]string{"add", "--db", dbPath, "--repeat=day", "--repeat-mode=schedule", "--repeat-start=2025-01-02", "--json", "Repeat me"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	var result repeatResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v: %q", err, out.String())
	}
	if !result.Verified || result.IDs.Created != "CREATED" || !result.Repeat.Active {
		t.Fatalf("unexpected result: %#v", result)
	}
}

type projectInsertingLauncher struct {
	dbPath string
	calls  int
}

func (l *projectInsertingLauncher) Open(args ...string) error {
	l.calls++
	if len(args) == 0 || !strings.HasPrefix(args[len(args)-1], "things:///") {
		return nil
	}
	conn, err := sql.Open("sqlite", l.dbPath)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Exec(`INSERT INTO TMTask (uuid,type,status,trashed,title,start,startBucket,creationDate,userModificationDate) VALUES ('CREATEDPROJ',1,0,0,'Repeat Project',1,4,?,?)`, float64(time.Now().Unix()), float64(time.Now().Unix()))
	return err
}

func TestRepeatAddProjectJSONDryRunDescribesAllStagesWithoutLaunch(t *testing.T) {
	dbPath := writeTestDB(t)
	launcher := &recordLauncher{}
	out := &bytes.Buffer{}
	app := &App{In: strings.NewReader(""), Out: out, Err: &bytes.Buffer{}, Launcher: launcher, DryRun: true}
	root := NewRoot(app)
	root.SetArgs([]string{"--dry-run", "add-project", "--db", dbPath, "--repeat=week", "--repeat-mode=schedule", "--repeat-start=2025-01-02", "--json", "Repeat Project"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if len(launcher.args) != 0 {
		t.Fatalf("dry-run launched Things: %#v", launcher.args)
	}
	var result repeatResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("stdout is not one JSON object: %v: %q", err, out.String())
	}
	if !result.DryRun || result.Repeat.Mode != "schedule" || result.Repeat.Unit != "week" || len(result.Stages) != 4 {
		t.Fatalf("unexpected preview: %#v", result)
	}
	if result.Database.Path == "" || result.Database.Source != "flag" || result.Intent == nil || !strings.HasPrefix(result.Intent.URL, "things:///add-project?") {
		t.Fatalf("missing preflight provenance/intent: %#v", result)
	}
}

func TestRepeatAddProjectAppliesAndVerifiesCreatedTemplate(t *testing.T) {
	dbPath := writeTestDB(t)
	launcher := &projectInsertingLauncher{dbPath: dbPath}
	out := &bytes.Buffer{}
	app := &App{In: strings.NewReader(""), Out: out, Err: &bytes.Buffer{}, Launcher: launcher}
	root := NewRoot(app)
	root.SetArgs([]string{"add-project", "--db", dbPath, "--repeat=week", "--repeat-mode=schedule", "--repeat-start=2025-01-02", "--json", "Repeat Project"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	var result repeatResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v: %q", err, out.String())
	}
	if !result.Verified || result.IDs.Created != "CREATEDPROJ" || !result.Repeat.Active {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestRepeatAddProjectWithCountOmitsEndDate(t *testing.T) {
	dbPath := writeTestDB(t)
	launcher := &projectInsertingLauncher{dbPath: dbPath}
	out := &bytes.Buffer{}
	app := &App{In: strings.NewReader(""), Out: out, Err: &bytes.Buffer{}, Launcher: launcher}
	root := NewRoot(app)
	root.SetArgs([]string{"add-project", "--db", dbPath, "--repeat=week", "--repeat-mode=schedule", "--repeat-start=2025-01-02", "--repeat-count=20", "--json", "Repeat Project"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	var result repeatResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v: %q", err, out.String())
	}
	if !result.Verified || result.Repeat.Count == nil || *result.Repeat.Count != 20 {
		t.Fatalf("expected count 20: %#v", result.Repeat)
	}
	if result.Repeat.EndDate != "" {
		t.Fatalf("count-based rule must not report an end date: %#v", result.Repeat)
	}
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	var rule []byte
	if err := conn.QueryRow(`SELECT rt1_recurrenceRule FROM TMTask WHERE uuid = 'CREATEDPROJ'`).Scan(&rule); err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if _, err := plist.Unmarshal(rule, &decoded); err != nil {
		t.Fatal(err)
	}
	if _, exists := decoded["ed"]; exists {
		t.Fatalf("count rule must omit ed: %#v", decoded)
	}
}

func TestRepeatAddProjectRejectsWhenBeforeLaunch(t *testing.T) {
	launcher := &recordLauncher{}
	app := &App{In: strings.NewReader(""), Out: &bytes.Buffer{}, Err: &bytes.Buffer{}, Launcher: launcher}
	root := NewRoot(app)
	root.SetArgs([]string{"add-project", "--repeat=day", "--when=today", "Nope"})
	if err := root.Execute(); err == nil || !strings.Contains(err.Error(), "cannot be combined") {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(launcher.args) != 0 {
		t.Fatal("launched Things")
	}
}

func TestRepeatAddProjectInvalidSpecFailsBeforeLaunch(t *testing.T) {
	launcher := &recordLauncher{}
	app := &App{In: strings.NewReader(""), Out: &bytes.Buffer{}, Err: &bytes.Buffer{}, Launcher: launcher}
	root := NewRoot(app)
	root.SetArgs([]string{"add-project", "--dry-run", "--repeat=day", "--repeat-every=0", "Nope"})
	if err := root.Execute(); err == nil {
		t.Fatal("expected invalid interval")
	}
	if len(launcher.args) != 0 {
		t.Fatalf("invalid spec launched Things: %#v", launcher.args)
	}
}

func TestRepeatAddProjectRejectsTerminalStatus(t *testing.T) {
	for _, arg := range []string{"--canceled", "--completed"} {
		launcher := &recordLauncher{}
		app := &App{In: strings.NewReader(""), Out: &bytes.Buffer{}, Err: &bytes.Buffer{}, Launcher: launcher}
		root := NewRoot(app)
		root.SetArgs([]string{"add-project", arg, "--repeat=day", "Nope"})
		if err := root.Execute(); err == nil || !strings.Contains(err.Error(), "cannot be combined") {
			t.Fatalf("unexpected error for %s: %v", arg, err)
		}
		if len(launcher.args) != 0 {
			t.Fatalf("%s launched Things: %#v", arg, launcher.args)
		}
	}
}

type recordLauncher struct {
	args []string
}

func (r *recordLauncher) Open(args ...string) error {
	r.args = append([]string{}, args...)
	return nil
}

func requireOpenURL(t *testing.T, launcher *recordLauncher) string {
	t.Helper()
	if len(launcher.args) != 2 {
		t.Fatalf("expected 2 args, got %d", len(launcher.args))
	}
	if launcher.args[0] != "-g" {
		t.Fatalf("expected -g flag, got %q", launcher.args[0])
	}
	if launcher.args[1] == "" {
		t.Fatalf("expected url arg, got empty")
	}
	return launcher.args[1]
}

func TestAddCommandWithTitle(t *testing.T) {
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"add", "New Todo"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	url := requireOpenURL(t, launcher)
	if !strings.Contains(url, "title=New%20Todo") {
		t.Fatalf("expected title in url, got %q", url)
	}
}

func TestAddCommandRejectsUnsafeTitle(t *testing.T) {
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"add", "tag=work"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err == nil {
		t.Fatalf("expected error")
	}
	if len(launcher.args) != 0 {
		t.Fatalf("expected no open invocation")
	}
}

func TestAddCommandAllowsUnsafeTitleWithFlag(t *testing.T) {
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"add", "--allow-unsafe-title", "tag=work"})
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

func TestAddCommandReadsStdin(t *testing.T) {
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader("Title\n\nNotes"),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"add", "-"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	url := requireOpenURL(t, launcher)
	if !strings.Contains(url, "title=Title") {
		t.Fatalf("expected title in url, got %q", url)
	}
	if !strings.Contains(url, "notes=Notes") {
		t.Fatalf("expected notes in url, got %q", url)
	}
}

func TestAddCommandShowQuickEntryWhenNoTitle(t *testing.T) {
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"add"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	url := requireOpenURL(t, launcher)
	if !strings.Contains(url, "show-quick-entry=true") {
		t.Fatalf("expected show-quick-entry in url, got %q", url)
	}
}
