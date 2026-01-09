package cli

import (
	"fmt"
	"strings"
)

var unsafeTitleSuggestions = map[string]string{
	"tag":      "--tags",
	"tags":     "--tags",
	"add-tags": "--add-tags",
}

func guardUnsafeTitle(title string, allowUnsafe bool) error {
	if allowUnsafe {
		return nil
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return nil
	}
	key := unsafeTitleKey(title)
	if key == "" {
		return nil
	}
	if suggestion, ok := unsafeTitleSuggestions[key]; ok {
		return fmt.Errorf("Error: title %q looks like %s=...; did you mean %s? Use --allow-unsafe-title to keep it as the title.", title, key, suggestion)
	}
	return fmt.Errorf("Error: title %q looks like %s=...; use --allow-unsafe-title to keep it as the title.", title, key)
}

func unsafeTitleKey(title string) string {
	parts := strings.Fields(strings.TrimSpace(title))
	if len(parts) == 0 {
		return ""
	}
	token := parts[0]
	token = strings.TrimPrefix(token, "--")
	token = strings.TrimPrefix(token, "-")
	idx := strings.Index(token, "=")
	if idx <= 0 {
		return ""
	}
	key := strings.ToLower(strings.TrimSpace(token[:idx]))
	if key == "" {
		return ""
	}
	switch key {
	case "tag", "tags", "add-tags", "when", "deadline", "list", "list-id", "area", "area-id", "project", "project-id", "heading", "notes", "id", "auth-token":
		return key
	}
	if strings.HasPrefix(key, "repeat-") {
		return key
	}
	return ""
}

func resolveWhenValue(value string, later bool) string {
	value = strings.TrimSpace(value)
	if value != "" {
		return value
	}
	if later {
		return "evening"
	}
	return ""
}

func validateWhenInput(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	switch strings.ToLower(value) {
	case "today", "tomorrow", "evening", "someday", "anytime", "inbox":
		return nil
	default:
		if _, _, err := parseDateOrTime(value); err != nil {
			msg := strings.TrimPrefix(err.Error(), "Error: ")
			return fmt.Errorf("Error: invalid --when value %q (%s)", value, msg)
		}
	}
	return nil
}
