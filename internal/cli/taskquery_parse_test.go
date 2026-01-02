package cli

import (
	"testing"

	"github.com/ossianhempel/things3-cli/internal/db"
)

func TestParseRichQueryMatch(t *testing.T) {
	expr, err := parseRichQuery("title:alpha AND tag:work")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tasks := []db.Task{
		{Title: "alpha task", Tags: []string{"work"}},
		{Title: "alpha task", Tags: []string{"home"}},
		{Title: "beta task", Tags: []string{"work"}},
	}
	filtered := filterTasksByQuery(tasks, expr)
	if len(filtered) != 1 {
		t.Fatalf("expected 1 match, got %d", len(filtered))
	}
	if filtered[0].Title != "alpha task" || filtered[0].Tags[0] != "work" {
		t.Fatalf("unexpected match: %+v", filtered[0])
	}
}

func TestParseRichQueryRegex(t *testing.T) {
	expr, err := parseRichQuery("title:/^alph/i")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tasks := []db.Task{
		{Title: "Alpha"},
		{Title: "beta"},
	}
	filtered := filterTasksByQuery(tasks, expr)
	if len(filtered) != 1 || filtered[0].Title != "Alpha" {
		t.Fatalf("unexpected matches: %+v", filtered)
	}
}

func TestParseRichQueryURLPredicate(t *testing.T) {
	expr, err := parseRichQuery("url:true")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tasks := []db.Task{
		{Notes: "see https://example.com"},
		{Notes: "no links here"},
	}
	filtered := filterTasksByQuery(tasks, expr)
	if len(filtered) != 1 {
		t.Fatalf("expected 1 match, got %d", len(filtered))
	}
}
