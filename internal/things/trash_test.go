package things

import "testing"

func TestBuildTrashScriptRequiresIDs(t *testing.T) {
	if _, err := BuildTrashScript(nil); err == nil {
		t.Fatalf("expected error")
	}
	if _, err := BuildTrashScript([]string{""}); err == nil {
		t.Fatalf("expected error")
	}
}

func TestBuildTrashScriptWithIDs(t *testing.T) {
	script, err := BuildTrashScript([]string{"123", "456"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(script, "delete to do id todoID") {
		t.Fatalf("expected delete in script, got %q", script)
	}
	if !contains(script, "\"123\"") || !contains(script, "\"456\"") {
		t.Fatalf("expected ids in script, got %q", script)
	}
}
