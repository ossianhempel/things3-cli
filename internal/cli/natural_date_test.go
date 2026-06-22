package cli

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestNormalizeWhenInputNaturalPhrases(t *testing.T) {
	withLocalTime(t, "Europe/Stockholm", func(now time.Time) {
		tests := []struct {
			input string
			want  string
		}{
			{"today", "2026-06-22"},
			{"tomorrow", "2026-06-23"},
			{"next Friday", "2026-06-26"},
			{"in 2 weeks", "2026-07-06"},
			{"in 1 month", "2026-07-22"},
			{"next week", "2026-06-29"},
		}

		for _, tt := range tests {
			got, err := normalizeWhenInput(tt.input, now)
			if err != nil {
				t.Fatalf("normalizeWhenInput(%q) failed: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("normalizeWhenInput(%q) = %q, want %q", tt.input, got, tt.want)
			}
		}
	})
}

func TestNormalizeWhenInputPreservesExplicitFormatsAndSpecialValues(t *testing.T) {
	withLocalTime(t, "Europe/Stockholm", func(now time.Time) {
		tests := []string{
			"2026-07-01",
			"2026-07-01 14:30",
			"2026-07-01 14:30:00",
			"2026-07-01T14:30:00+02:00",
			"evening",
			"anytime",
			"someday",
			"inbox",
		}

		for _, input := range tests {
			got, err := normalizeWhenInput(input, now)
			if err != nil {
				t.Fatalf("normalizeWhenInput(%q) failed: %v", input, err)
			}
			if got != input {
				t.Fatalf("normalizeWhenInput(%q) = %q, want original", input, got)
			}
		}
	})
}

func TestNormalizeScheduleInputRejectsInvalidNaturalDate(t *testing.T) {
	withLocalTime(t, "Europe/Stockholm", func(now time.Time) {
		for _, input := range []string{"next someday", "in two weeks", "in 0 days", "this Friday"} {
			if _, err := normalizeWhenInput(input, now); err == nil {
				t.Fatalf("expected normalizeWhenInput(%q) to fail", input)
			}
		}
	})
}

func TestNormalizeDeadlineInputTimezoneLocalDate(t *testing.T) {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		t.Fatalf("load location: %v", err)
	}
	previous := time.Local
	time.Local = loc
	t.Cleanup(func() { time.Local = previous })

	now := time.Date(2026, 6, 22, 6, 30, 0, 0, time.UTC)
	got, err := normalizeDeadlineInput("tomorrow", now)
	if err != nil {
		t.Fatalf("normalizeDeadlineInput failed: %v", err)
	}
	if got != "2026-06-22" {
		t.Fatalf("got %q, want local tomorrow 2026-06-22", got)
	}
}

func TestUpdateCommandDryRunShowsNormalizedNaturalDates(t *testing.T) {
	withLocalTime(t, "Europe/Stockholm", func(_ time.Time) {
		out := &bytes.Buffer{}
		app := &App{
			In:       strings.NewReader(""),
			Out:      out,
			Err:      &bytes.Buffer{},
			Launcher: &recordLauncher{},
		}

		root := NewRoot(app)
		root.SetArgs([]string{
			"--dry-run",
			"update",
			"--auth-token", "tok",
			"--id", "123",
			"--when", "next Friday",
			"--deadline", "in 2 weeks",
			"--no-verify",
		})
		root.SetOut(app.Out)
		root.SetErr(app.Err)

		if err := root.Execute(); err != nil {
			t.Fatalf("execute failed: %v", err)
		}

		url := out.String()
		for _, want := range []string{"when=2026-06-26", "deadline=2026-07-06"} {
			if !strings.Contains(url, want) {
				t.Fatalf("expected %q in dry-run URL, got %q", want, url)
			}
		}
		if strings.Contains(url, "next%20Friday") || strings.Contains(url, "in%202%20weeks") {
			t.Fatalf("expected normalized URL, got %q", url)
		}
	})
}

func TestAddCommandRejectsInvalidDeadline(t *testing.T) {
	withLocalTime(t, "Europe/Stockholm", func(_ time.Time) {
		launcher := &recordLauncher{}
		app := &App{
			In:       strings.NewReader(""),
			Out:      &bytes.Buffer{},
			Err:      &bytes.Buffer{},
			Launcher: launcher,
		}

		root := NewRoot(app)
		root.SetArgs([]string{"add", "--deadline", "next someday", "Title"})
		root.SetOut(app.Out)
		root.SetErr(app.Err)

		if err := root.Execute(); err == nil {
			t.Fatalf("expected error")
		}
		if len(launcher.args) != 0 {
			t.Fatalf("expected no open invocation")
		}
	})
}

func withLocalTime(t *testing.T, name string, fn func(time.Time)) {
	t.Helper()
	loc, err := time.LoadLocation(name)
	if err != nil {
		t.Fatalf("load location: %v", err)
	}
	previous := time.Local
	time.Local = loc
	t.Cleanup(func() { time.Local = previous })

	fn(time.Date(2026, 6, 22, 9, 30, 0, 0, loc))
}
