package cli

import (
	"bytes"
	"strings"
	"testing"
)

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
