package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootHelpCommand(t *testing.T) {
	launcher := &recordLauncher{}
	out := &bytes.Buffer{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      out,
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"help"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err != nil && err != ErrHelpPrinted {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "things - manage todos with Things 3") {
		t.Fatalf("unexpected help output: %q", out.String())
	}
	if len(launcher.args) != 0 {
		t.Fatalf("expected no open invocation")
	}
}

func TestAddHelpCommand(t *testing.T) {
	launcher := &recordLauncher{}
	out := &bytes.Buffer{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      out,
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"help", "add"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err != nil && err != ErrHelpPrinted {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "things add - add new todo") {
		t.Fatalf("unexpected help output: %q", out.String())
	}
	for _, want := range []string{"after-completion mode (the default)", "fixed calendar cadence", "--repeat-start", "--dry-run", "template state"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("add help missing %q: %q", want, out.String())
		}
	}
}

func TestAddProjectHelpRejectsRepeatingProjects(t *testing.T) {
	out := &bytes.Buffer{}
	app := &App{In: strings.NewReader(""), Out: out, Err: &bytes.Buffer{}, Launcher: &recordLauncher{}}
	root := NewRoot(app)
	root.SetArgs([]string{"help", "add-project"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)
	if err := root.Execute(); err != nil && err != ErrHelpPrinted {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "Repeating projects are not supported") {
		t.Fatalf("add-project help should state the limitation: %q", out.String())
	}
	for _, unsupported := range []string{"--repeat=UNIT", "--repeat-mode=MODE", "--repeat-start=DATE"} {
		if strings.Contains(out.String(), unsupported) {
			t.Fatalf("add-project help advertises unsupported option %q", unsupported)
		}
	}
}

func TestUpdateHelpExplainsRepeatPermissionsAndRecovery(t *testing.T) {
	out := &bytes.Buffer{}
	app := &App{In: strings.NewReader(""), Out: out, Err: &bytes.Buffer{}, Launcher: &recordLauncher{}}
	root := NewRoot(app)
	root.SetArgs([]string{"help", "update"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)
	if err := root.Execute(); err != nil && err != ErrHelpPrinted {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{"Repeat-only updates", "Full Disk Access", "require both", "partial result", "trusted UUID", "--repeat-clear"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("update help missing %q: %q", want, out.String())
		}
	}
}

func TestInvalidHelpCommand(t *testing.T) {
	launcher := &recordLauncher{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"help", "nope"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err == nil {
		t.Fatalf("expected error")
	}
}

func TestRootHelpFlag(t *testing.T) {
	launcher := &recordLauncher{}
	out := &bytes.Buffer{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      out,
		Err:      &bytes.Buffer{},
		Launcher: launcher,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"-h"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err != nil && err != ErrHelpPrinted {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "COMMANDS") {
		t.Fatalf("unexpected help output: %q", out.String())
	}
}
