package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestDeleteCommandRequiresTarget(t *testing.T) {
	runner := &recordScriptRunner{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Scripter: runner,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"delete"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err == nil {
		t.Fatalf("expected error")
	}
	if runner.script != "" {
		t.Fatalf("expected no script execution")
	}
}

func TestDeleteCommandWithID(t *testing.T) {
	runner := &recordScriptRunner{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Scripter: runner,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"delete", "--id", "ABC123", "--confirm", "ABC123"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	script := requireScript(t, runner)
	if !strings.Contains(script, "first to do whose id is \"ABC123\"") {
		t.Fatalf("expected id lookup in script, got %q", script)
	}
	if !strings.Contains(script, "delete targetTodo") {
		t.Fatalf("expected delete in script, got %q", script)
	}
}

func TestDeleteCommandWithTitle(t *testing.T) {
	runner := &recordScriptRunner{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Scripter: runner,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"delete", "Pay bills", "--confirm", "Pay bills"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	script := requireScript(t, runner)
	if !strings.Contains(script, "to do \"Pay bills\"") {
		t.Fatalf("expected title lookup in script, got %q", script)
	}
	if !strings.Contains(script, "delete targetTodo") {
		t.Fatalf("expected delete in script, got %q", script)
	}
}

func TestDeleteCommandQueryRequiresConfirmation(t *testing.T) {
	dbPath := writeTestDB(t)
	runner := &recordScriptRunner{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Scripter: runner,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"delete", "--db", dbPath, "--search", "Task One"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err == nil {
		t.Fatalf("expected error")
	}
	if runner.script != "" {
		t.Fatalf("expected no script execution")
	}
}

func TestDeleteCommandQueryWithYes(t *testing.T) {
	dbPath := writeTestDB(t)
	runner := &recordScriptRunner{}
	app := &App{
		In:       strings.NewReader(""),
		Out:      &bytes.Buffer{},
		Err:      &bytes.Buffer{},
		Scripter: runner,
	}

	root := NewRoot(app)
	root.SetArgs([]string{"delete", "--db", dbPath, "--search", "Task One", "--yes"})
	root.SetOut(app.Out)
	root.SetErr(app.Err)

	if err := root.Execute(); err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	script := requireScript(t, runner)
	if !strings.Contains(script, "delete to do id") {
		t.Fatalf("expected trash script, got %q", script)
	}
}
