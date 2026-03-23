package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func writeActionLog(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "actions.jsonl")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write action log: %v", err)
	}
	return path
}

func readActionLog(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read action log: %v", err)
	}
	return string(data)
}

func TestRemoveLastActionFromPath(t *testing.T) {
	t.Run("removes last line", func(t *testing.T) {
		path := writeActionLog(t, `{"type":"update","items":[]}`+"\n"+`{"type":"trash","items":[]}`+"\n")
		if err := removeLastActionFromPath(path); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := readActionLog(t, path)
		want := `{"type":"update","items":[]}` + "\n"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("removes only entry", func(t *testing.T) {
		path := writeActionLog(t, `{"type":"update","items":[]}`+"\n")
		if err := removeLastActionFromPath(path); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := readActionLog(t, path)
		if got != "" {
			t.Errorf("expected empty file, got %q", got)
		}
	})

	t.Run("handles trailing newlines", func(t *testing.T) {
		path := writeActionLog(t, `{"type":"update","items":[]}`+"\n"+`{"type":"trash","items":[]}`+"\n\n\n")
		if err := removeLastActionFromPath(path); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := readActionLog(t, path)
		want := `{"type":"update","items":[]}` + "\n\n\n"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("returns error for missing file", func(t *testing.T) {
		err := removeLastActionFromPath(filepath.Join(t.TempDir(), "nonexistent.jsonl"))
		if err == nil {
			t.Error("expected error for missing file, got nil")
		}
	})
}
