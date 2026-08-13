package repeat

import (
	"testing"
	"time"

	"howett.net/plist"
)

func TestBuildUpdateWeeklySchedule(t *testing.T) {
	anchor := time.Date(2026, 1, 6, 15, 0, 0, 0, time.Local)
	spec := Spec{
		Mode:   ModeSchedule,
		Unit:   UnitWeek,
		Every:  2,
		Anchor: anchor,
	}
	update, err := BuildUpdate(spec)
	if err != nil {
		t.Fatalf("BuildUpdate failed: %v", err)
	}

	var decoded map[string]any
	if _, err := plist.Unmarshal(update.RecurrenceRule, &decoded); err != nil {
		t.Fatalf("unmarshal rule: %v", err)
	}
	assertInt(t, decoded["fa"], 2)
	assertInt(t, decoded["fu"], 256)
	assertInt(t, decoded["tp"], 0)
	assertInt(t, decoded["ts"], 0)

	offsets := decoded["of"].([]any)
	if len(offsets) != 1 {
		t.Fatalf("expected 1 offset, got %d", len(offsets))
	}
	offset := offsets[0].(map[string]any)
	assertInt(t, offset["wd"], int(anchor.Weekday()))

	expectedStart := thingsDateValue(anchor.AddDate(0, 0, 1))
	if update.InstanceCreationStartDate != expectedStart {
		t.Fatalf("start date mismatch: got %d want %d", update.InstanceCreationStartDate, expectedStart)
	}
	expectedNext := thingsDateValue(anchor.AddDate(0, 0, 14))
	if update.NextInstanceStartDate == nil || *update.NextInstanceStartDate != expectedNext {
		t.Fatalf("expected next instance for scheduled date")
	}
	if update.SetDeadline {
		t.Fatalf("unexpected deadline flag")
	}
}

func TestBuildUpdateDeadlineOffset(t *testing.T) {
	offset := 3
	spec := Spec{
		Mode:           ModeAfterCompletion,
		Unit:           UnitMonth,
		Every:          1,
		Anchor:         time.Date(2026, 1, 15, 9, 0, 0, 0, time.Local),
		DeadlineOffset: &offset,
	}
	update, err := BuildUpdate(spec)
	if err != nil {
		t.Fatalf("BuildUpdate failed: %v", err)
	}
	if !update.SetDeadline {
		t.Fatalf("expected deadline flag")
	}
	if update.Deadline == nil {
		t.Fatalf("expected deadline sentinel")
	}
	var decoded map[string]any
	if _, err := plist.Unmarshal(update.RecurrenceRule, &decoded); err != nil {
		t.Fatalf("unmarshal rule: %v", err)
	}
	assertInt(t, decoded["ts"], -3)
}

func TestBuildUpdateUntilDate(t *testing.T) {
	end := time.Date(2026, 2, 1, 0, 0, 0, 0, time.Local)
	spec := Spec{
		Mode:    ModeSchedule,
		Unit:    UnitDay,
		Every:   1,
		Anchor:  time.Date(2026, 1, 6, 0, 0, 0, 0, time.Local),
		EndDate: &end,
	}
	update, err := BuildUpdate(spec)
	if err != nil {
		t.Fatalf("BuildUpdate failed: %v", err)
	}
	var decoded map[string]any
	if _, err := plist.Unmarshal(update.RecurrenceRule, &decoded); err != nil {
		t.Fatalf("unmarshal rule: %v", err)
	}
	assertInt(t, decoded["ed"], int(end.Unix()))
}

func TestBuildUpdateRepeatCountOmitsEndDate(t *testing.T) {
	count := 20
	spec := Spec{
		Mode:   ModeSchedule,
		Unit:   UnitWeek,
		Every:  1,
		Anchor: time.Date(2026, 1, 5, 0, 0, 0, 0, time.Local),
		Count:  &count,
	}
	update, err := BuildUpdate(spec)
	if err != nil {
		t.Fatalf("BuildUpdate failed: %v", err)
	}
	var decoded map[string]any
	if _, err := plist.Unmarshal(update.RecurrenceRule, &decoded); err != nil {
		t.Fatalf("unmarshal rule: %v", err)
	}
	assertInt(t, decoded["rc"], 20)
	if _, exists := decoded["ed"]; exists {
		t.Fatalf("count-based rule must omit ed")
	}
	if update.NextInstanceStartDate == nil {
		t.Fatalf("expected a next instance for schedule mode")
	}
}

func TestBuildUpdateDeadlineShiftsIntervalAnchor(t *testing.T) {
	// A Monday anchor with a 6-day deadline must keep occurrences on Monday:
	// ia = anchor + deadlineOffset, ts = -deadlineOffset, wd = weekday(ia).
	offset := 6
	anchor := time.Date(2026, 8, 17, 0, 0, 0, 0, time.Local) // Monday
	spec := Spec{
		Mode:           ModeSchedule,
		Unit:           UnitWeek,
		Every:          1,
		Anchor:         anchor,
		DeadlineOffset: &offset,
	}
	update, err := BuildUpdate(spec)
	if err != nil {
		t.Fatalf("BuildUpdate failed: %v", err)
	}
	var decoded map[string]any
	if _, err := plist.Unmarshal(update.RecurrenceRule, &decoded); err != nil {
		t.Fatalf("unmarshal rule: %v", err)
	}
	assertInt(t, decoded["ts"], -6)
	iaSec := int64(decoded["ia"].(float64))
	ia := time.Unix(iaSec, 0)
	if ia.Weekday() != anchor.AddDate(0, 0, 6).Weekday() {
		t.Fatalf("ia should shift by deadline offset, got %v", ia)
	}
	offsets := decoded["of"].([]any)
	offsetDict := offsets[0].(map[string]any)
	assertInt(t, offsetDict["wd"], int(ia.Weekday()))
	if !update.SetDeadline || update.Deadline == nil {
		t.Fatalf("expected deadline sentinel")
	}
}

func TestBuildUpdateFutureAnchorNextIsAnchorItself(t *testing.T) {
	anchor := time.Now().AddDate(0, 0, 4) // future
	spec := Spec{
		Mode:   ModeSchedule,
		Unit:   UnitWeek,
		Every:  1,
		Anchor: anchor,
	}
	update, err := BuildUpdate(spec)
	if err != nil {
		t.Fatalf("BuildUpdate failed: %v", err)
	}
	if update.NextInstanceStartDate == nil {
		t.Fatalf("expected next instance")
	}
	expected := thingsDateValue(normalizeDate(anchor))
	if *update.NextInstanceStartDate != expected {
		t.Fatalf("future anchor should be its own first instance, got %d want %d", *update.NextInstanceStartDate, expected)
	}
}

func TestBuildUpdateDeadlineShiftsEndBoundary(t *testing.T) {
	// A bounded schedule with a deadline must keep its encoded end boundary in
	// the same shifted frame as ia so the final occurrence is not dropped.
	offset := 6
	anchor := time.Date(2026, 8, 17, 0, 0, 0, 0, time.Local) // Monday
	end := time.Date(2026, 8, 17, 0, 0, 0, 0, time.Local)     // first occurrence
	spec := Spec{
		Mode:           ModeSchedule,
		Unit:           UnitWeek,
		Every:          1,
		Anchor:         anchor,
		EndDate:        &end,
		DeadlineOffset: &offset,
	}
	update, err := BuildUpdate(spec)
	if err != nil {
		t.Fatalf("BuildUpdate failed: %v", err)
	}
	var decoded map[string]any
	if _, err := plist.Unmarshal(update.RecurrenceRule, &decoded); err != nil {
		t.Fatalf("unmarshal rule: %v", err)
	}
	iaSec := int64(decoded["ia"].(float64))
	edSec := int64(decoded["ed"].(float64))
	if edSec < iaSec {
		t.Fatalf("encoded end boundary must stay at or after the shifted anchor: ia=%d ed=%d", iaSec, edSec)
	}
	if got := time.Unix(edSec, 0).Format("2006-01-02"); got != "2026-08-23" {
		t.Fatalf("encoded end date = %s, want 2026-08-23 (anchor + deadline offset)", got)
	}
}

func assertInt(t *testing.T, value any, expected int) {
	t.Helper()
	switch v := value.(type) {
	case int:
		if v != expected {
			t.Fatalf("expected %d got %d", expected, v)
		}
	case int64:
		if int(v) != expected {
			t.Fatalf("expected %d got %d", expected, v)
		}
	case uint64:
		if int(v) != expected {
			t.Fatalf("expected %d got %d", expected, v)
		}
	case float64:
		if int(v) != expected {
			t.Fatalf("expected %d got %d", expected, int(v))
		}
	default:
		t.Fatalf("unexpected int type %T", value)
	}
}
