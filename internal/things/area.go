package things

import (
	"fmt"
	"strings"
)

// AddAreaOptions defines options for add-area.
type AddAreaOptions struct {
	Tags string
}

// UpdateAreaOptions defines options for update-area.
type UpdateAreaOptions struct {
	ID      string
	Tags    string
	AddTags string
}

// BuildAddAreaScript builds an AppleScript snippet for creating an area.
func BuildAddAreaScript(opts AddAreaOptions, rawInput string) (string, error) {
	title := parseSingleLineTitle(rawInput)
	if title == "" {
		return "", errMissingTitle
	}

	var b strings.Builder
	b.WriteString("tell application \"Things3\"\n")
	b.WriteString("  set newArea to make new area with properties {name:\"")
	b.WriteString(escapeAppleScriptString(title))
	b.WriteString("\"}\n")
	if strings.TrimSpace(opts.Tags) != "" {
		b.WriteString("  set tag names of newArea to \"")
		b.WriteString(escapeAppleScriptString(opts.Tags))
		b.WriteString("\"\n")
	}
	b.WriteString("end tell")
	return b.String(), nil
}

// BuildUpdateAreaScript builds an AppleScript snippet for updating an area.
func BuildUpdateAreaScript(opts UpdateAreaOptions, rawInput string) (string, error) {
	title := parseSingleLineTitle(rawInput)
	if opts.ID == "" && title == "" {
		return "", errMissingAreaTarget
	}
	if strings.TrimSpace(opts.Tags) == "" && strings.TrimSpace(opts.AddTags) == "" {
		return "", errMissingAreaTags
	}

	target := areaTarget(opts.ID, title)
	var b strings.Builder
	b.WriteString("tell application \"Things3\"\n")
	b.WriteString("  set targetArea to ")
	b.WriteString(target)
	b.WriteString("\n")

	if strings.TrimSpace(opts.Tags) != "" {
		b.WriteString("  set tag names of targetArea to \"")
		b.WriteString(escapeAppleScriptString(opts.Tags))
		b.WriteString("\"\n")
	} else {
		addTags := strings.TrimSpace(opts.AddTags)
		b.WriteString("  set currentTags to tag names of targetArea\n")
		b.WriteString("  if currentTags is missing value then set currentTags to \"\"\n")
		b.WriteString("  if currentTags is \"\" then\n")
		b.WriteString("    set tag names of targetArea to \"")
		b.WriteString(escapeAppleScriptString(addTags))
		b.WriteString("\"\n")
		b.WriteString("  else\n")
		b.WriteString("    set tag names of targetArea to currentTags & \", \" & \"")
		b.WriteString(escapeAppleScriptString(addTags))
		b.WriteString("\"\n")
		b.WriteString("  end if\n")
	}

	b.WriteString("end tell")
	return b.String(), nil
}

func parseSingleLineTitle(rawInput string) string {
	if rawInput == "" {
		return ""
	}
	title := rawInput
	if HasMultipleLines(rawInput) {
		title = FindTitle(rawInput)
	}
	return strings.TrimSpace(title)
}

func areaTarget(id string, title string) string {
	if id == "" {
		return fmt.Sprintf("area \"%s\"", escapeAppleScriptString(title))
	}
	return fmt.Sprintf("first area whose id is \"%s\"", escapeAppleScriptString(id))
}
