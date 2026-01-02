package cli

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ossianhempel/things3-cli/internal/db"
)

type ActionType string

const (
	ActionUpdate ActionType = "update"
	ActionTrash  ActionType = "trash"
)

type ActionEntry struct {
	Timestamp string       `json:"timestamp"`
	Type      ActionType   `json:"type"`
	Items     []ActionItem `json:"items"`
}

type ActionItem struct {
	UUID         string   `json:"uuid"`
	Title        string   `json:"title"`
	Status       int      `json:"status"`
	Notes        string   `json:"notes,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	Deadline     string   `json:"deadline,omitempty"`
	Start        string   `json:"start,omitempty"`
	StartDate    string   `json:"start_date,omitempty"`
	ProjectID    string   `json:"project_id,omitempty"`
	AreaID       string   `json:"area_id,omitempty"`
	HeadingTitle string   `json:"heading_title,omitempty"`
}

func actionLogPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, "things3-cli")
	if err := os.MkdirAll(path, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(path, "actions.jsonl"), nil
}

func appendAction(entry ActionEntry) error {
	path, err := actionLogPath()
	if err != nil {
		return err
	}
	entry.Timestamp = time.Now().Format(time.RFC3339)
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	return enc.Encode(entry)
}

func readLastAction() (ActionEntry, error) {
	path, err := actionLogPath()
	if err != nil {
		return ActionEntry{}, err
	}
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ActionEntry{}, errors.New("no actions logged")
		}
		return ActionEntry{}, err
	}
	defer file.Close()

	var lastLine string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lastLine = line
		}
	}
	if err := scanner.Err(); err != nil {
		return ActionEntry{}, err
	}
	if lastLine == "" {
		return ActionEntry{}, errors.New("no actions logged")
	}
	var entry ActionEntry
	if err := json.Unmarshal([]byte(lastLine), &entry); err != nil {
		return ActionEntry{}, err
	}
	return entry, nil
}

func removeLastAction() error {
	path, err := actionLogPath()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.TrimSpace(lines[i]) != "" {
			lines = append(lines[:i], lines[i+1:]...)
			break
		}
	}
	content := strings.Join(lines, "\n")
	return os.WriteFile(path, []byte(content), 0o644)
}

func taskToActionItem(task db.Task) ActionItem {
	return ActionItem{
		UUID:         task.UUID,
		Title:        task.Title,
		Status:       task.Status,
		Notes:        task.Notes,
		Tags:         task.Tags,
		Deadline:     task.Deadline,
		Start:        task.Start,
		StartDate:    task.StartDate,
		ProjectID:    task.ProjectID,
		AreaID:       task.AreaID,
		HeadingTitle: task.HeadingTitle,
	}
}
