package things

import "strings"

func escapeAppleScriptString(input string) string {
	replaced := strings.ReplaceAll(input, "\\", "\\\\")
	return strings.ReplaceAll(replaced, "\"", "\\\"")
}
