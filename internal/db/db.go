package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// Store wraps a Things database connection.
type Store struct {
	conn *sql.DB
	path string
}

// Open opens a Things database at the provided path in read-only mode.
func Open(path string) (*Store, error) {
	if path == "" {
		return nil, fmt.Errorf("empty database path")
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve database path: %w", err)
	}
	dsn := sqliteDSN(abs, "ro")
	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if err := conn.Ping(); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("open database: %w", err)
	}
	return &Store{conn: conn, path: abs}, nil
}

// OpenDefault resolves the Things database path and opens it.
func OpenDefault(override string) (*Store, string, error) {
	path, err := ResolveDatabasePath(override)
	if err != nil {
		return nil, "", err
	}
	store, err := Open(path)
	if err != nil {
		return nil, path, err
	}
	return store, path, nil
}

// Close closes the underlying database connection.
func (s *Store) Close() error {
	if s == nil || s.conn == nil {
		return nil
	}
	return s.conn.Close()
}

// Path returns the resolved database path.
func (s *Store) Path() string {
	if s == nil {
		return ""
	}
	return s.path
}

// OpenWritable opens a Things database at the provided path in read-write mode.
func OpenWritable(path string) (*Store, error) {
	if path == "" {
		return nil, fmt.Errorf("empty database path")
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve database path: %w", err)
	}
	info, err := os.Lstat(abs)
	if err != nil {
		return nil, fmt.Errorf("inspect writable database: %w", err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("unsafe writable database: symlinks are not allowed")
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("unsafe writable database: target is not a regular file")
	}
	canonical, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return nil, fmt.Errorf("canonicalize writable database: %w", err)
	}
	dsn := sqliteDSN(canonical, "rw")
	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if err := conn.Ping(); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("open database: %w", err)
	}
	openedInfo, err := os.Stat(canonical)
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("reinspect writable database: %w", err)
	}
	if !os.SameFile(info, openedInfo) {
		_ = conn.Close()
		return nil, fmt.Errorf("unsafe writable database: target changed while opening")
	}
	var taskTable string
	if err := conn.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='TMTask'`).Scan(&taskTable); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("unrecognized Things database: TMTask table not found")
	}
	return &Store{conn: conn, path: canonical}, nil
}

// OpenDefaultWritable resolves the Things database path and opens it in read-write mode.
func OpenDefaultWritable(override string) (*Store, string, error) {
	path, err := ResolveDatabasePath(override)
	if err != nil {
		return nil, "", err
	}
	store, err := OpenWritable(path)
	if err != nil {
		return nil, path, err
	}
	return store, path, nil
}

func sqliteDSN(path string, mode string) string {
	u := url.URL{Scheme: "file", Path: path}
	q := u.Query()
	if mode == "" {
		mode = "ro"
	}
	q.Set("mode", mode)
	u.RawQuery = q.Encode()
	return u.String()
}
