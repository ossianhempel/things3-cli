package cli

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ossianhempel/things3-cli/internal/db"
)

const whenVerifyTimeout = 4 * time.Second

func openVerifyStore(app *App, dbPath string) (*db.Store, error) {
	store, _, err := db.OpenDefault(dbPath)
	if err != nil {
		if app != nil {
			msg := formatDBError(err).Error()
			msg = strings.TrimPrefix(msg, "Error: ")
			fmt.Fprintf(app.Err, "Warning: could not verify update (Things database unavailable): %s\n", msg)
		}
		return nil, nil
	}
	return store, nil
}

func verifyWhenApplied(store *db.Store, id string, expected string) error {
	expected = strings.TrimSpace(expected)
	if store == nil || expected == "" {
		return nil
	}
	deadline := time.Now().Add(whenVerifyTimeout)
	var lastTask *db.Task
	for time.Now().Before(deadline) {
		task, err := store.TaskByID(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return err
		}
		lastTask = task
		if whenMatches(*task, expected) {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	if lastTask == nil {
		return fmt.Errorf("Error: failed to verify update for %s", id)
	}
	if lastTask.Repeating {
		return fmt.Errorf("Error: cannot update when for repeating todos (id %s)", id)
	}
	return fmt.Errorf("Error: update did not apply (expected when=%s, got start=%q start_date=%q). Check THINGS_AUTH_TOKEN and Things permissions.", expected, lastTask.Start, lastTask.StartDate)
}

func whenMatches(task db.Task, expected string) bool {
	expected = strings.TrimSpace(expected)
	if expected == "" {
		return true
	}
	switch strings.ToLower(expected) {
	case "inbox":
		return strings.EqualFold(task.Start, "Inbox")
	case "anytime":
		return strings.EqualFold(task.Start, "Anytime") && task.StartDate == ""
	case "someday":
		return strings.EqualFold(task.Start, "Someday") && task.StartDate == ""
	case "today", "evening":
		return task.StartDate == dateString(time.Now())
	case "tomorrow":
		return task.StartDate == dateString(time.Now().AddDate(0, 0, 1))
	default:
		parsed, _, err := parseDateOrTime(expected)
		if err != nil {
			return true
		}
		return task.StartDate == dateString(parsed)
	}
}

func dateString(t time.Time) string {
	return t.In(time.Local).Format("2006-01-02")
}

func validateEveningTask(task db.Task, allowNonToday bool) error {
	if allowNonToday {
		return nil
	}
	if task.StartDate == "" {
		return nil
	}
	today := dateString(time.Now())
	if task.StartDate != today {
		return fmt.Errorf("Error: refusing to move task %s to This Evening because it is scheduled for %s (use --allow-non-today to override)", task.UUID, task.StartDate)
	}
	return nil
}
