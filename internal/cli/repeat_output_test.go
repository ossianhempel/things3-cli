package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/ossianhempel/things3-cli/internal/db"
	"github.com/ossianhempel/things3-cli/internal/repeat"
)

func TestRenderRepeatResultJSONIsSingleStableObject(t *testing.T) {
	offset := 2
	result := repeatResult{
		SchemaVersion: 1, Action: "apply", IDs: repeatIDs{Requested: "REQ", Template: "TPL"},
		Repeat: &db.RepeatState{Mode: "schedule", Unit: "day", Interval: 1, Anchor: "2025-01-02", DeadlineOffset: &offset, Active: true, Scheduled: true},
		Stages: []repeatStage{{Name: "database", Status: "applied"}, {Name: "verification", Status: "verified"}}, Verified: true,
	}
	var out bytes.Buffer
	if err := renderRepeatResult(&out, result, true); err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	dec := json.NewDecoder(&out)
	if err := dec.Decode(&decoded); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if dec.More() {
		t.Fatal("expected exactly one JSON object")
	}
	if decoded["schema_version"] != float64(1) || decoded["action"] != "apply" || decoded["verified"] != true {
		t.Fatalf("unexpected contract: %#v", decoded)
	}
	serialized, _ := json.Marshal(decoded)
	for _, forbidden := range []string{"recurrenceRule", "plist", "notes"} {
		if strings.Contains(string(serialized), forbidden) {
			t.Fatalf("leaked %s: %s", forbidden, serialized)
		}
	}
}

func TestRecoveryPreservesSemanticsAndShellQuotes(t *testing.T) {
	offset := 2
	until, _ := time.ParseInLocation("2006-01-02", "2025-04-01", time.Local)
	anchor, _ := time.ParseInLocation("2006-01-02", "2025-01-02", time.Local)
	spec := RepeatSpec{Spec: repeat.Spec{Mode: repeat.ModeSchedule, Unit: repeat.UnitWeek, Every: 3, Anchor: anchor, EndDate: &until, DeadlineOffset: &offset}}
	argv := recoveryArgv("odd ' uuid", spec, "/tmp/other database.sqlite")
	joined := strings.Join(argv, " ")
	for _, want := range []string{"--db /tmp/other database.sqlite", "--repeat week", "--repeat-mode schedule", "--repeat-every 3", "--repeat-until 2025-04-01", "--repeat-deadline 2"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("missing %q in %#v", want, argv)
		}
	}
	if got := shellRender(argv); !strings.Contains(got, `'odd '\'' uuid'`) {
		t.Fatalf("unsafe render: %s", got)
	}
}

func TestRenderRepeatResultHumanNamesLifecycleAndUUIDs(t *testing.T) {
	var out bytes.Buffer
	err := renderRepeatResult(&out, repeatResult{SchemaVersion: 1, Action: "clear", IDs: repeatIDs{Requested: "REQ", Template: "TPL"}, Stages: []repeatStage{{Name: "database", Status: "cleared"}}, Verified: true}, false)
	if err != nil {
		t.Fatal(err)
	}
	text := out.String()
	for _, want := range []string{"Repeat clear: verified", "Requested UUID: REQ", "Template UUID: TPL", "Stage database: cleared"} {
		if !strings.Contains(text, want) {
			t.Fatalf("missing %q in %q", want, text)
		}
	}
}
