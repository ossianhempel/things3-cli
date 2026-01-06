package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/ossianhempel/things3-cli/internal/db"
	"github.com/ossianhempel/things3-cli/internal/repeat"
	"github.com/ossianhempel/things3-cli/internal/things"
)

func extractTitle(rawInput string, titlesRaw string) string {
	if strings.TrimSpace(titlesRaw) != "" {
		return ""
	}
	rawInput = strings.TrimSpace(rawInput)
	if rawInput == "" {
		return ""
	}
	if things.HasMultipleLines(rawInput) {
		return strings.TrimSpace(things.FindTitle(rawInput))
	}
	return rawInput
}

func waitForCreatedItem(store *db.Store, title string, taskType int, started time.Time) (string, error) {
	if store == nil {
		return "", fmt.Errorf("database not initialized")
	}
	if strings.TrimSpace(title) == "" {
		return "", fmt.Errorf("title required to locate created item")
	}
	deadline := time.Now().Add(90 * time.Second)
	since := float64(started.Unix())
	for time.Now().Before(deadline) {
		matches, err := store.TasksByTitleSince(title, taskType, since)
		if err != nil {
			return "", err
		}
		if len(matches) == 1 {
			return matches[0].UUID, nil
		}
		if len(matches) > 1 {
			return "", fmt.Errorf("multiple items created with title %q; use --id", title)
		}
		time.Sleep(200 * time.Millisecond)
	}
	return "", fmt.Errorf("timed out waiting for the created item")
}

func resolveRepeatTarget(store *db.Store, id string, expectedType int) (string, bool, error) {
	target, err := store.RepeatTargetByID(id)
	if err != nil {
		return "", false, err
	}
	resolvedID := id
	usedTemplate := false
	if target.RepeatingTemplateID != "" {
		resolvedID = target.RepeatingTemplateID
		usedTemplate = true
		target, err = store.RepeatTargetByID(resolvedID)
		if err != nil {
			return "", usedTemplate, err
		}
	}
	if target.Type != expectedType {
		return "", usedTemplate, fmt.Errorf("Error: item type mismatch for repeat update")
	}
	if target.Trashed {
		return "", usedTemplate, fmt.Errorf("Error: cannot update repeating rules for trashed items")
	}
	if target.Status != db.StatusIncomplete {
		return "", usedTemplate, fmt.Errorf("Error: repeating rules require an incomplete item")
	}
	return resolvedID, usedTemplate, nil
}

func applyRepeatSpec(store *db.Store, id string, spec RepeatSpec) error {
	if spec.Clear {
		return store.ClearRepeatRule(id)
	}
	update, err := repeat.BuildUpdate(spec.Spec)
	if err != nil {
		return err
	}
	return store.ApplyRepeatRule(id, update)
}
