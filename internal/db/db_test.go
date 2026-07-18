package db

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOpenWritableRejectsSymlink(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "main.sqlite")
	if err := os.WriteFile(target, []byte("not sqlite"), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "linked.sqlite")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	if _, err := OpenWritable(link); err == nil || !strings.Contains(err.Error(), "symlinks are not allowed") {
		t.Fatalf("expected symlink rejection, got %v", err)
	}
}

func TestOpenWritableRejectsUnrecognizedDatabase(t *testing.T) {
	path := filepath.Join(t.TempDir(), "other.sqlite")
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(`CREATE TABLE Other (id INTEGER)`); err != nil {
		t.Fatal(err)
	}
	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := OpenWritable(path); err == nil || !strings.Contains(err.Error(), "TMTask table not found") {
		t.Fatalf("expected provenance rejection, got %v", err)
	}
}
