package things

import (
	"fmt"
	"strings"
)

// BuildTrashScript builds an AppleScript snippet to move todos to Trash.
func BuildTrashScript(ids []string) (string, error) {
	if len(ids) == 0 {
		return "", fmt.Errorf("Error: Must specify --id=ID or query")
	}
	quoted := make([]string, 0, len(ids))
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		quoted = append(quoted, "\""+escapeAppleScriptString(id)+"\"")
	}
	if len(quoted) == 0 {
		return "", fmt.Errorf("Error: Must specify --id=ID or query")
	}

	var b strings.Builder
	b.WriteString("tell application \"Things3\"\n")
	b.WriteString("  repeat with todoID in {")
	b.WriteString(strings.Join(quoted, ", "))
	b.WriteString("}\n")
	b.WriteString("    try\n")
	b.WriteString("      delete to do id todoID\n")
	b.WriteString("    end try\n")
	b.WriteString("  end repeat\n")
	b.WriteString("end tell")
	return b.String(), nil
}
