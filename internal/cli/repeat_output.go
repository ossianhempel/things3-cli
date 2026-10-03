package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ossianhempel/things3-cli/internal/db"
	"github.com/ossianhempel/things3-cli/internal/repeat"
)

type repeatStage struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type repeatStageName string
type repeatStageStatus string

const (
	repeatStageURL          repeatStageName = "url"
	repeatStageLocate       repeatStageName = "locate"
	repeatStageDatabase     repeatStageName = "database"
	repeatStageVerification repeatStageName = "verification"

	repeatStatusPlanned   repeatStageStatus = "planned"
	repeatStatusCompleted repeatStageStatus = "completed"
	repeatStatusFailed    repeatStageStatus = "failed"
	repeatStatusApplied   repeatStageStatus = "applied"
	repeatStatusCleared   repeatStageStatus = "cleared"
	repeatStatusVerified  repeatStageStatus = "verified"
)

type repeatIDs struct {
	Requested string `json:"requested,omitempty"`
	Created   string `json:"created,omitempty"`
	Template  string `json:"template,omitempty"`
}

type repeatDatabase struct {
	Path   string `json:"path"`
	Source string `json:"source"`
}

type repeatIntent struct {
	URL string `json:"url,omitempty"`
}

type repeatRecovery struct {
	Argv []string `json:"argv"`
}

type repeatResult struct {
	SchemaVersion int              `json:"schema_version"`
	Action        string           `json:"action"`
	DryRun        bool             `json:"dry_run"`
	Partial       bool             `json:"partial"`
	IDs           repeatIDs        `json:"ids"`
	Repeat        *db.RepeatState  `json:"repeat,omitempty"`
	Stages        []repeatStage    `json:"stages"`
	Verified      bool             `json:"verified"`
	Database      repeatDatabase   `json:"database"`
	Intent        *repeatIntent    `json:"intent,omitempty"`
	Recovery      []repeatRecovery `json:"recovery,omitempty"`
}

type repeatUpdatePlan struct {
	store        *db.Store
	targetID     string
	usedTemplate bool
	clear        bool
	taskType     int
	prepared     db.RepeatUpdate
	expected     *db.RepeatState
	spec         RepeatSpec
	result       repeatResult
}

func prepareRepeatUpdate(dbPath, requestedID string, spec RepeatSpec, dryRun bool) (*repeatUpdatePlan, error) {
	return prepareRepeatUpdateForType(dbPath, requestedID, spec, dryRun, db.TaskTypeTodo)
}

func prepareRepeatUpdateForType(dbPath, requestedID string, spec RepeatSpec, dryRun bool, taskType int) (*repeatUpdatePlan, error) {
	var store *db.Store
	var err error
	var resolvedPath string
	if dryRun {
		store, resolvedPath, err = db.OpenDefault(dbPath)
	} else {
		store, resolvedPath, err = db.OpenDefaultWritable(dbPath)
	}
	if err != nil {
		return nil, err
	}
	fail := func(err error) (*repeatUpdatePlan, error) {
		_ = store.Close()
		return nil, err
	}
	if err := store.ValidateRepeatSchema(); err != nil {
		return fail(err)
	}
	targetID, usedTemplate, err := resolveRepeatTarget(store, requestedID, taskType)
	if err != nil {
		return fail(err)
	}
	target, err := store.RepeatTargetByID(targetID)
	if err != nil {
		return fail(err)
	}
	plan := &repeatUpdatePlan{
		store:        store,
		targetID:     targetID,
		usedTemplate: usedTemplate,
		clear:        spec.Clear,
		taskType:     taskType,
		spec:         spec,
		result: repeatResult{
			SchemaVersion: 1,
			Action:        "apply",
			DryRun:        dryRun,
			IDs:           repeatIDs{Requested: requestedID, Template: targetID},
			Database:      repeatDatabase{Path: store.Path(), Source: repeatDatabaseSource(dbPath, resolvedPath)},
		},
	}
	if spec.Spec.Count != nil && target.InstanceCount >= *spec.Spec.Count {
		// A repeat count is a total for the template. On an existing template
		// that has already spawned instances, the requested total must be
		// larger than the number already created, or the schedule would either
		// stop immediately or drop prior occurrences.
		return fail(fmt.Errorf("Error: --repeat-count=%d is not greater than the %d occurrences already created; pass a total that includes prior occurrences or use --repeat-clear first", *spec.Spec.Count, target.InstanceCount))
	}
	if spec.Clear {
		plan.result.Action = "clear"
		return plan, nil
	}
	if spec.Spec.Anchor.IsZero() {
		spec.Spec.Anchor = time.Now()
		plan.spec = spec
	}
	plan.prepared, err = repeat.BuildUpdate(spec.Spec)
	if err != nil {
		return fail(err)
	}
	plan.expected = expectedRepeatState(spec.Spec)
	plan.result.Repeat = plan.expected
	return plan, nil
}

func (p *repeatUpdatePlan) executeMutation() error {
	var err error
	if p.clear {
		err = p.store.ClearRepeatRule(p.targetID)
	} else {
		err = p.store.ApplyRepeatRule(p.targetID, p.prepared)
	}
	if err != nil {
		p.result.Partial = len(p.result.Stages) > 0
		p.result.addStage(repeatStageDatabase, repeatStatusFailed)
		p.result.Recovery = []repeatRecovery{{Argv: recoveryArgvForType(p.targetID, p.spec, p.result.Database.Path, p.taskType)}}
		return err
	}
	status := repeatStatusApplied
	if p.clear {
		status = repeatStatusCleared
	}
	p.result.addStage(repeatStageDatabase, status)
	actual, err := verifyRepeatState(p.store, p.targetID, p.expected, p.clear)
	p.result.Repeat = actual
	if err != nil {
		p.result.failStage(repeatStageVerification)
		p.result.Recovery = []repeatRecovery{{Argv: recoveryArgvForType(p.targetID, p.spec, p.result.Database.Path, p.taskType)}}
		return err
	}
	p.result.markVerified()
	return nil
}

func (r *repeatResult) addStage(name repeatStageName, status repeatStageStatus) {
	r.Stages = append(r.Stages, repeatStage{Name: string(name), Status: string(status)})
}

func (r *repeatResult) failStage(name repeatStageName) {
	r.Partial = true
	r.addStage(name, repeatStatusFailed)
}

func (r *repeatResult) markVerified() {
	r.Verified = true
	r.addStage(repeatStageVerification, repeatStatusVerified)
}

func renderRepeatResult(out io.Writer, result repeatResult, asJSON bool) error {
	if asJSON {
		return json.NewEncoder(out).Encode(result)
	}
	status := "planned"
	if result.Verified {
		status = "verified"
	}
	if result.Partial {
		status = "partial"
	}
	fmt.Fprintf(out, "Repeat %s: %s\n", result.Action, status)
	if result.IDs.Requested != "" {
		fmt.Fprintf(out, "Requested UUID: %s\n", result.IDs.Requested)
	}
	if result.IDs.Created != "" {
		fmt.Fprintf(out, "Created UUID: %s\n", result.IDs.Created)
	}
	if result.IDs.Template != "" {
		fmt.Fprintf(out, "Template UUID: %s\n", result.IDs.Template)
	}
	if result.Database.Path != "" {
		fmt.Fprintf(out, "Database: %s (%s)\n", result.Database.Path, result.Database.Source)
	}
	if result.Intent != nil && result.Intent.URL != "" {
		fmt.Fprintf(out, "URL intent: %s\n", result.Intent.URL)
	}
	if result.Repeat != nil {
		fmt.Fprintf(out, "Rule: %s every %d %s; anchor %s", result.Repeat.Mode, result.Repeat.Interval, result.Repeat.Unit, result.Repeat.Anchor)
		if result.Repeat.EndDate != "" {
			fmt.Fprintf(out, "; until %s", result.Repeat.EndDate)
		}
		if result.Repeat.Count != nil {
			fmt.Fprintf(out, "; count %d", *result.Repeat.Count)
		}
		if result.Repeat.DeadlineOffset != nil {
			fmt.Fprintf(out, "; deadline offset %d", *result.Repeat.DeadlineOffset)
		}
		fmt.Fprintln(out)
		fmt.Fprintf(out, "State: active=%t paused=%t scheduled=%t\n", result.Repeat.Active, result.Repeat.Paused, result.Repeat.Scheduled)
	}
	for _, stage := range result.Stages {
		fmt.Fprintf(out, "Stage %s: %s\n", stage.Name, stage.Status)
	}
	for _, recovery := range result.Recovery {
		fmt.Fprintf(out, "Recovery: %s\n", shellRender(recovery.Argv))
	}
	return nil
}

func expectedRepeatState(spec repeat.Spec) *db.RepeatState {
	mode := "after-completion"
	if spec.Mode == repeat.ModeSchedule {
		mode = "schedule"
	}
	unit := map[repeat.Unit]string{repeat.UnitDay: "day", repeat.UnitWeek: "week", repeat.UnitMonth: "month", repeat.UnitYear: "year"}[spec.Unit]
	anchor := spec.Anchor
	state := &db.RepeatState{Active: true, Paused: false, Scheduled: true, Mode: mode, Unit: unit, Interval: spec.Every, Anchor: anchor.Format("2006-01-02"), DeadlineOffset: spec.DeadlineOffset, Count: spec.Count}
	if spec.EndDate != nil {
		state.EndDate = spec.EndDate.Format("2006-01-02")
	}
	return state
}

func repeatDatabaseSource(override, resolved string) string {
	if strings.TrimSpace(override) != "" {
		return "flag"
	}
	if strings.TrimSpace(os.Getenv("THINGSDB")) != "" {
		return "environment"
	}
	if resolved != "" {
		return "discovered"
	}
	return "unknown"
}

func redactURLIntent(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return "things:///update?<redacted>"
	}
	q := u.Query()
	if q.Has("auth-token") {
		q.Set("auth-token", "REDACTED")
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func recoveryArgv(id string, spec RepeatSpec, dbPath string) []string {
	return recoveryArgvForType(id, spec, dbPath, db.TaskTypeTodo)
}

func recoveryArgvForType(id string, spec RepeatSpec, dbPath string, taskType int) []string {
	command := "update"
	if taskType == db.TaskTypeProject {
		command = "update-project"
	}
	argv := []string{"things", command, "--id", id}
	if strings.TrimSpace(dbPath) != "" {
		argv = append(argv, "--db", dbPath)
	}
	if spec.Clear {
		return append(argv, "--repeat-clear")
	}
	unit := map[repeat.Unit]string{repeat.UnitDay: "day", repeat.UnitWeek: "week", repeat.UnitMonth: "month", repeat.UnitYear: "year"}[spec.Spec.Unit]
	mode := "after-completion"
	if spec.Spec.Mode == repeat.ModeSchedule {
		mode = "schedule"
	}
	argv = append(argv, "--repeat", unit, "--repeat-mode", mode, "--repeat-every", fmt.Sprint(spec.Spec.Every), "--repeat-start", spec.Spec.Anchor.Format("2006-01-02"))
	if spec.Spec.EndDate != nil {
		argv = append(argv, "--repeat-until", spec.Spec.EndDate.Format("2006-01-02"))
	}
	if spec.Spec.DeadlineOffset != nil {
		argv = append(argv, "--repeat-deadline", fmt.Sprint(*spec.Spec.DeadlineOffset))
	}
	if spec.Spec.Count != nil {
		argv = append(argv, "--repeat-count", fmt.Sprint(*spec.Spec.Count))
	}
	return argv
}

func shellRender(argv []string) string {
	quoted := make([]string, len(argv))
	for i, arg := range argv {
		quoted[i] = "'" + strings.ReplaceAll(arg, "'", "'\\''") + "'"
	}
	return strings.Join(quoted, " ")
}

func verifyRepeatState(store *db.Store, id string, expected *db.RepeatState, clearing bool) (*db.RepeatState, error) {
	actual, err := store.RepeatStateByID(id)
	if err != nil {
		return nil, err
	}
	if clearing {
		if actual.Active || actual.Mode != "" || actual.Unit != "" {
			return actual, fmt.Errorf("repeat clear verification failed for %s", id)
		}
		return actual, nil
	}
	if !actual.Active || !actual.Scheduled || actual.Paused || actual.DecodeWarning != "" || actual.Mode != expected.Mode || actual.Unit != expected.Unit || actual.Interval != expected.Interval || actual.Anchor != expected.Anchor || actual.EndDate != expected.EndDate {
		return actual, fmt.Errorf("repeat verification failed for %s", id)
	}
	if (actual.DeadlineOffset == nil) != (expected.DeadlineOffset == nil) || actual.DeadlineOffset != nil && *actual.DeadlineOffset != *expected.DeadlineOffset {
		return actual, fmt.Errorf("repeat deadline verification failed for %s", id)
	}
	if (actual.Count == nil) != (expected.Count == nil) || actual.Count != nil && *actual.Count != *expected.Count {
		return actual, fmt.Errorf("repeat count verification failed for %s", id)
	}
	return actual, nil
}
