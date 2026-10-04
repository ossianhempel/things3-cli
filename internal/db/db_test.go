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

func TestOpenWritableRejectsThingsManagedDatabase(t *testing.T) {
	for _, parent := range []string{"Things Database.thingsdatabase", "JLMPQHK86H.com.culturedcode.ThingsMac"} {
		t.Run(parent, func(t *testing.T) {
			dir := filepath.Join(t.TempDir(), parent)
			if err := os.MkdirAll(dir, 0o700); err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(dir, "main.sqlite")
			original := []byte("must not be opened or changed")
			if err := os.WriteFile(path, original, 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := OpenWritable(path); err == nil || !strings.Contains(err.Error(), "live Things database are unsupported") {
				t.Fatalf("expected native repeat guidance, got %v", err)
			}
			got, err := os.ReadFile(path)
			if err != nil || string(got) != string(original) {
				t.Fatal("database changed")
			}
			link := filepath.Join(t.TempDir(), "alias")
			if err := os.Symlink(dir, link); err != nil {
				t.Fatal(err)
			}
			if _, err := OpenWritable(filepath.Join(link, "main.sqlite")); err == nil || !strings.Contains(err.Error(), "live Things database are unsupported") {
				t.Fatalf("parent symlink bypassed guard: %v", err)
			}
		})
	}
}
