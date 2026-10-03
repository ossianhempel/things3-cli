package cli

import (
	"fmt"
	"time"

	"github.com/ossianhempel/things3-cli/internal/repeat"
	"github.com/spf13/cobra"
)

// RepeatOptions captures CLI repeat flags before parsing.
type RepeatOptions struct {
	Rule           string
	Mode           string
	Every          int
	Start          string
	Until          string
	DeadlineOffset int
	Count          int
	Clear          bool
}

// RepeatSpec wraps a parsed repeat spec.
type RepeatSpec struct {
	Enabled bool
	Clear   bool
	Spec    repeat.Spec
}

func addRepeatFlags(cmd *cobra.Command, opts *RepeatOptions, allowClear bool) {
	flags := cmd.Flags()
	flags.StringVar(&opts.Rule, "repeat", "", "Repeat unit: day, week, month, or year")
	flags.StringVar(&opts.Mode, "repeat-mode", "after-completion", "Repeat mode: after-completion or schedule")
	flags.IntVar(&opts.Every, "repeat-every", 1, "Repeat interval (every N units)")
	flags.StringVar(&opts.Start, "repeat-start", "", "Repeat anchor date (YYYY-MM-DD)")
	flags.StringVar(&opts.Until, "repeat-until", "", "Repeat until date (YYYY-MM-DD)")
	flags.IntVar(&opts.DeadlineOffset, "repeat-deadline", 0, "Add repeating deadlines (days earlier)")
	flags.IntVar(&opts.Count, "repeat-count", 0, "Stop after this many occurrences (omit for unlimited)")
	if allowClear {
		flags.BoolVar(&opts.Clear, "repeat-clear", false, "Remove repeating schedule")
	}
}

func parseRepeatSpec(cmd *cobra.Command, opts RepeatOptions) (RepeatSpec, error) {
	if opts.Clear {
		if opts.Rule != "" ||
			cmd.Flags().Changed("repeat-mode") ||
			cmd.Flags().Changed("repeat-every") ||
			cmd.Flags().Changed("repeat-start") ||
			cmd.Flags().Changed("repeat-until") ||
			cmd.Flags().Changed("repeat-deadline") ||
			cmd.Flags().Changed("repeat-count") {
			return RepeatSpec{}, fmt.Errorf("Error: --repeat-clear cannot be combined with other repeat flags")
		}
		return RepeatSpec{Enabled: true, Clear: true}, nil
	}

	changed := opts.Rule != "" ||
		cmd.Flags().Changed("repeat-mode") ||
		cmd.Flags().Changed("repeat-every") ||
		cmd.Flags().Changed("repeat-start") ||
		cmd.Flags().Changed("repeat-until") ||
		cmd.Flags().Changed("repeat-deadline") ||
		cmd.Flags().Changed("repeat-count")

	if !changed {
		return RepeatSpec{Enabled: false}, nil
	}
	if opts.Rule == "" {
		return RepeatSpec{}, fmt.Errorf("Error: --repeat is required when using repeat flags")
	}
	if cmd.Flags().Changed("repeat-count") && opts.Count <= 0 {
		return RepeatSpec{}, fmt.Errorf("Error: --repeat-count must be a positive integer")
	}
	if cmd.Flags().Changed("repeat-count") && cmd.Flags().Changed("repeat-until") {
		return RepeatSpec{}, fmt.Errorf("Error: --repeat-count and --repeat-until are mutually exclusive; use one to bound the schedule")
	}

	mode, err := repeat.ParseMode(opts.Mode)
	if err != nil {
		return RepeatSpec{}, fmt.Errorf("Error: %v", err)
	}
	unit, err := repeat.ParseUnit(opts.Rule)
	if err != nil {
		return RepeatSpec{}, fmt.Errorf("Error: %v", err)
	}
	anchor := time.Now()
	if opts.Start != "" {
		parsed, _, err := parseDateOrTime(opts.Start)
		if err != nil {
			return RepeatSpec{}, err
		}
		anchor = parsed
	}
	var until *time.Time
	if opts.Until != "" {
		parsed, _, err := parseDateOrTime(opts.Until)
		if err != nil {
			return RepeatSpec{}, err
		}
		until = &parsed
	}

	var deadlineOffset *int
	if cmd.Flags().Changed("repeat-deadline") {
		value := opts.DeadlineOffset
		deadlineOffset = &value
	}

	var count *int
	if cmd.Flags().Changed("repeat-count") {
		value := opts.Count
		count = &value
	}

	spec := repeat.Spec{
		Mode:           mode,
		Unit:           unit,
		Every:          opts.Every,
		Anchor:         anchor,
		EndDate:        until,
		DeadlineOffset: deadlineOffset,
		Count:          count,
	}

	return RepeatSpec{Enabled: true, Spec: spec}, nil
}
