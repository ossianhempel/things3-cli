package cli

import (
	"bytes"
	"encoding/json"
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
	startBucket := 1
	referenceDate := 135004288
	task := db.Task{UUID: "ABC", Title: "Task", Status: db.StatusCompleted, StartBucket: &startBucket, TodayIndexReferenceDate: &referenceDate}
	opts := TaskOutputOptions{
		Format: "csv",
		Select: []string{"uuid", "title", "status", "start_bucket", "today_index_reference_date"},
	}
	var buf bytes.Buffer
	if err := writeTasks(&buf, []db.Task{task, {UUID: "NULL", Title: "Null", Status: db.StatusIncomplete}}, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := strings.TrimSpace(buf.String())
	lines := strings.Split(out, "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "UUID,TITLE,STATUS,START_BUCKET,TODAY_INDEX_REFERENCE_DATE" {
		t.Fatalf("unexpected header: %q", lines[0])
	}
	if !strings.Contains(lines[1], "ABC,Task,completed,1,135004288") {
		t.Fatalf("unexpected row: %q", lines[1])
	}
	if lines[2] != "NULL,Null,incomplete,," {
		t.Fatalf("unexpected nil row: %q", lines[2])
	}
}

func TestTodayIndexReferenceDateAliases(t *testing.T) {
	for _, alias := range []string{"today_index_reference_date", "todayindexreferencedate", "today-index-reference-date"} {
		if got := normalizeTaskField(alias); got != "today_index_reference_date" {
			t.Errorf("normalizeTaskField(%q) = %q", alias, got)
		}
	}
}

func TestTodayIndexReferenceDateStructuredNullBehavior(t *testing.T) {
	referenceDate := 135004288
	tasks := []db.Task{
		{UUID: "PRESENT", Title: "Present", TodayIndexReferenceDate: &referenceDate},
		{UUID: "NULL", Title: "Null"},
	}

	t.Run("full JSON omits nil", func(t *testing.T) {
		var buf bytes.Buffer
		if err := writeTasks(&buf, tasks, TaskOutputOptions{Format: "json"}); err != nil {
			t.Fatalf("write tasks: %v", err)
		}
		var records []map[string]any
		if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
			t.Fatalf("decode JSON: %v", err)
		}
		if got := records[0]["today_index_reference_date"]; got != float64(referenceDate) {
			t.Fatalf("present value = %#v", got)
		}
		if _, ok := records[1]["today_index_reference_date"]; ok {
			t.Fatalf("nil full JSON field should be omitted: %#v", records[1])
		}
	})

	t.Run("selected JSON preserves null", func(t *testing.T) {
		var buf bytes.Buffer
		if err := writeTasks(&buf, tasks, TaskOutputOptions{Format: "json", Select: []string{"uuid", "today_index_reference_date"}}); err != nil {
			t.Fatalf("write tasks: %v", err)
		}
		var records []map[string]any
		if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
			t.Fatalf("decode JSON: %v", err)
		}
		if got := records[0]["today_index_reference_date"]; got != float64(referenceDate) {
			t.Fatalf("present selected value = %#v", got)
		}
		if got, ok := records[1]["today_index_reference_date"]; !ok || got != nil {
			t.Fatalf("nil selected field = %#v, present=%v", got, ok)
		}
	})

	t.Run("full JSONL omits nil", func(t *testing.T) {
		var buf bytes.Buffer
		if err := writeTasks(&buf, tasks, TaskOutputOptions{Format: "jsonl"}); err != nil {
			t.Fatalf("write tasks: %v", err)
		}
		lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
		if !strings.Contains(lines[0], `"today_index_reference_date":135004288`) {
			t.Fatalf("present JSONL field missing: %s", lines[0])
		}
		if strings.Contains(lines[1], "today_index_reference_date") {
			t.Fatalf("nil full JSONL field should be omitted: %s", lines[1])
		}
	})

	t.Run("selected JSONL preserves null", func(t *testing.T) {
		var buf bytes.Buffer
		if err := writeTasks(&buf, tasks, TaskOutputOptions{Format: "jsonl", Select: []string{"uuid", "today_index_reference_date"}}); err != nil {
			t.Fatalf("write tasks: %v", err)
		}
		lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
		if !strings.Contains(lines[0], `"today_index_reference_date":135004288`) || !strings.Contains(lines[1], `"today_index_reference_date":null`) {
			t.Fatalf("unexpected selected JSONL: %v", lines)
		}
	})

	t.Run("selected table leaves nil blank", func(t *testing.T) {
		var buf bytes.Buffer
		if err := writeTasks(&buf, tasks[1:], TaskOutputOptions{Format: "table", Select: []string{"uuid", "today_index_reference_date"}}); err != nil {
			t.Fatalf("write tasks: %v", err)
		}
		if strings.Contains(buf.String(), "<nil>") {
			t.Fatalf("nil table field should be blank: %q", buf.String())
		}
	})
}

func TestRepeatSemanticFieldsAcrossSelectedFormats(t *testing.T) {
	offset := 3
	task := db.Task{UUID: "TPL", Repeat: &db.RepeatState{Active: true, Mode: "schedule", Unit: "week", Interval: 2, Anchor: "2025-01-02", DeadlineOffset: &offset, Scheduled: true}}
	fields, err := parseTaskSelect("uuid,repeat_active,repeat_mode,repeat_unit,repeat_interval,repeat_anchor,repeat_deadline_offset,repeat_scheduled")
	if err != nil {
		t.Fatal(err)
	}
	for _, format := range []string{"json", "jsonl", "csv", "table"} {
		t.Run(format, func(t *testing.T) {
			var out bytes.Buffer
			if err := writeTasks(&out, []db.Task{task}, TaskOutputOptions{Format: format, Select: fields}); err != nil {
				t.Fatal(err)
			}
			for _, want := range []string{"TPL", "schedule", "week", "2025-01-02"} {
				if !strings.Contains(out.String(), want) {
					t.Fatalf("%s missing %q: %q", format, want, out.String())
				}
			}
		})
	}
}
