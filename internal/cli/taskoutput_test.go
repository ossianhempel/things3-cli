package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ossianhempel/things3-cli/internal/db"
)

func TestResolveTaskOutputOptionsDefaults(t *testing.T) {
	opts, err := resolveTaskOutputOptions("", false, "", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Format != "table" {
		t.Fatalf("expected table format, got %q", opts.Format)
	}

	opts, err = resolveTaskOutputOptions("", true, "", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Format != "json" {
		t.Fatalf("expected json format, got %q", opts.Format)
	}
}

func TestResolveTaskOutputOptionsInvalidFormat(t *testing.T) {
	if _, err := resolveTaskOutputOptions("nope", false, "", false); err == nil {
		t.Fatalf("expected error for invalid format")
	}
	if _, err := resolveTaskOutputOptions("csv", true, "", false); err == nil {
		t.Fatalf("expected error for json + non-json format")
	}
}

func TestWriteTasksCSVSelect(t *testing.T) {
	task := db.Task{UUID: "ABC", Title: "Task", Status: db.StatusCompleted}
	opts := TaskOutputOptions{
		Format: "csv",
		Select: []string{"uuid", "title", "status"},
	}
	var buf bytes.Buffer
	if err := writeTasks(&buf, []db.Task{task}, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := strings.TrimSpace(buf.String())
	lines := strings.Split(out, "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0] != "UUID,TITLE,STATUS" {
		t.Fatalf("unexpected header: %q", lines[0])
	}
	if !strings.Contains(lines[1], "ABC,Task,completed") {
		t.Fatalf("unexpected row: %q", lines[1])
	}
}
