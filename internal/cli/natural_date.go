package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var weekdayNames = map[string]time.Weekday{
	"sunday":    time.Sunday,
	"sun":       time.Sunday,
	"monday":    time.Monday,
	"mon":       time.Monday,
	"tuesday":   time.Tuesday,
	"tue":       time.Tuesday,
	"tues":      time.Tuesday,
	"wednesday": time.Wednesday,
	"wed":       time.Wednesday,
	"thursday":  time.Thursday,
	"thu":       time.Thursday,
	"thur":      time.Thursday,
	"thurs":     time.Thursday,
	"friday":    time.Friday,
	"fri":       time.Friday,
	"saturday":  time.Saturday,
	"sat":       time.Saturday,
}

func normalizeWhenInput(value string, now time.Time) (string, error) {
	return normalizeScheduleInput("--when", value, now, true)
}

func normalizeDeadlineInput(value string, now time.Time) (string, error) {
	return normalizeScheduleInput("--deadline", value, now, false)
}

func normalizeScheduleInput(flag, value string, now time.Time, allowWhenSpecial bool) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", nil
	}

	lower := strings.ToLower(value)
	if allowWhenSpecial {
		switch lower {
		case "evening", "anytime", "someday", "inbox":
			return lower, nil
		}
	}

	_, _, explicitErr := parseDateOrTime(value)
	if explicitErr == nil {
		return value, nil
	}

	parsed, err := parseNaturalLocalDate(value, now)
	if err != nil {
		if looksExplicitDate(value) {
			msg := strings.TrimPrefix(explicitErr.Error(), "Error: ")
			return "", fmt.Errorf("Error: invalid %s value %q (%s)", flag, value, msg)
		}
		return "", fmt.Errorf("Error: invalid %s value %q (%s)", flag, value, err)
	}
	return parsed.Format("2006-01-02"), nil
}

func parseNaturalLocalDate(input string, now time.Time) (time.Time, error) {
	normalized := strings.Join(strings.Fields(strings.ToLower(strings.TrimSpace(input))), " ")
	today := localDate(now)

	switch normalized {
	case "today":
		return today, nil
	case "tomorrow":
		return today.AddDate(0, 0, 1), nil
	case "next week":
		return today.AddDate(0, 0, 7), nil
	case "next month":
		return today.AddDate(0, 1, 0), nil
	case "next year":
		return today.AddDate(1, 0, 0), nil
	}

	parts := strings.Fields(normalized)
	if len(parts) == 2 && parts[0] == "next" {
		weekday, ok := weekdayNames[parts[1]]
		if !ok {
			return time.Time{}, fmt.Errorf("unsupported natural date")
		}
		return nextWeekday(today, weekday), nil
	}

	if len(parts) == 3 && parts[0] == "in" {
		count, err := strconv.Atoi(parts[1])
		if err != nil || count < 1 {
			return time.Time{}, fmt.Errorf("relative dates require a positive whole number")
		}
		unit := strings.TrimSuffix(parts[2], "s")
		switch unit {
		case "day":
			return today.AddDate(0, 0, count), nil
		case "week":
			return today.AddDate(0, 0, count*7), nil
		case "month":
			return today.AddDate(0, count, 0), nil
		case "year":
			return today.AddDate(count, 0, 0), nil
		default:
			return time.Time{}, fmt.Errorf("unsupported relative date unit %q", parts[2])
		}
	}

	return time.Time{}, fmt.Errorf("unsupported natural date")
}

func localDate(t time.Time) time.Time {
	local := t.In(time.Local)
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, time.Local)
}

func nextWeekday(today time.Time, target time.Weekday) time.Time {
	days := (int(target) - int(today.Weekday()) + 7) % 7
	if days == 0 {
		days = 7
	}
	return today.AddDate(0, 0, days)
}

func looksExplicitDate(value string) bool {
	if strings.HasPrefix(strings.ToLower(strings.TrimSpace(value)), "in ") {
		return false
	}
	for _, r := range value {
		if r >= '0' && r <= '9' {
			return true
		}
	}
	return false
}
