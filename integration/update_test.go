package integration_test

import "testing"

func TestUpdateAuthTokenRequired(t *testing.T) {
	_, errOut, code := runThings(t, "", "update")
	requireFailure(t, code)
	assertContains(t, errOut, "Missing Things auth token")
}

func TestUpdateIDRequired(t *testing.T) {
	_, errOut, code := runThings(t, "", "update", "--auth-token=token")
	requireFailure(t, code)
	assertContains(t, errOut, "Must specify --id=id")
}

func TestUpdateWhenOption(t *testing.T) {
	out, _, code := runThings(t, "", "update", "--auth-token=token", "--id=1", "--when=2021-05-20")
	requireSuccess(t, code)
	assertContains(t, out, "when=2021-05-20")
}

func TestUpdateDeadlineOption(t *testing.T) {
	out, _, code := runThings(t, "", "update", "--auth-token=token", "--id=1", "--deadline=2021-05-20")
	requireSuccess(t, code)
	assertContains(t, out, "deadline=2021-05-20")
}

func TestUpdateCanceledPrecedence(t *testing.T) {
	out, _, code := runThings(t, "", "update", "--auth-token=token", "--id=1", "--canceled", "--completed")
	requireSuccess(t, code)
	assertContains(t, out, "canceled=true")
	assertNotContains(t, out, "completed=true")
}

func TestUpdateCompletedOption(t *testing.T) {
	out, _, code := runThings(t, "", "update", "--auth-token=token", "--id=1", "--completed")
	requireSuccess(t, code)
	assertContains(t, out, "completed=true")

	out, _, code = runThings(t, "", "update", "--auth-token=token", "--id=1", "--cancelled")
	requireSuccess(t, code)
	assertContains(t, out, "canceled=true")
}

func TestUpdateChecklistItemOption(t *testing.T) {
	out, _, code := runThings(t, "", "update", "--auth-token=token", "--id=1", "--checklist-item=ITEM 1", "--checklist-item=ITEM 2")
	requireSuccess(t, code)
	assertContains(t, out, join(enc("ITEM 1"), enc("ITEM 2")))
}

func TestUpdateCreationDateOption(t *testing.T) {
	out, _, code := runThings(t, "", "update", "--auth-token=token", "--id=1", "--creation-date=2021-05-20")
	requireSuccess(t, code)
	assertContains(t, out, "creation-date=2021-05-20")
}

func TestUpdateCompletionDateOption(t *testing.T) {
	out, _, code := runThings(t, "", "update", "--auth-token=token", "--id=1", "--completion-date=2021-05-20")
	requireSuccess(t, code)
	assertContains(t, out, "completion-date=2021-05-20")
}

func TestUpdateListOption(t *testing.T) {
	out, _, code := runThings(t, "", "update", "--auth-token=token", "--id=1", "--list=LIST")
	requireSuccess(t, code)
	assertContains(t, out, "list=LIST")
}

func TestUpdateListIDOption(t *testing.T) {
	out, _, code := runThings(t, "", "update", "--auth-token=token", "--id=1", "--list-id=1")
	requireSuccess(t, code)
	assertContains(t, out, "list-id=1")
}

func TestUpdateListIDPrecedence(t *testing.T) {
	out, _, code := runThings(t, "", "update", "--auth-token=token", "--id=1", "--list-id=1", "--list=LIST")
	requireSuccess(t, code)
	assertContains(t, out, "list-id=1")
	assertNotContains(t, out, "list=LIST")
}

func TestUpdateHeadingOption(t *testing.T) {
	out, _, code := runThings(t, "", "update", "--auth-token=token", "--id=1", "--heading=THE HEADING")
	requireSuccess(t, code)
	assertContains(t, out, enc("THE HEADING"))
}

func TestUpdateRevealOption(t *testing.T) {
	out, _, code := runThings(t, "", "update", "--auth-token=token", "--id=1", "--reveal")
	requireSuccess(t, code)
	assertContains(t, out, "reveal=true")
}

func TestUpdateDuplicateOption(t *testing.T) {
	out, _, code := runThings(t, "", "update", "--auth-token=token", "--id=1", "--duplicate")
	requireSuccess(t, code)
	assertContains(t, out, "duplicate=true")
}

func TestUpdateNotesOption(t *testing.T) {
	out, _, code := runThings(t, "", "update", "--auth-token=token", "--id=1", "--notes=NOTES NOTES")
	requireSuccess(t, code)
	assertContains(t, out, "notes="+enc("NOTES NOTES"))
}

func TestUpdateAppendNotesOption(t *testing.T) {
	out, _, code := runThings(t, "", "update", "--auth-token=token", "--id=1", "--append-notes=APPEND NOTES")
	requireSuccess(t, code)
	assertContains(t, out, "append-notes="+enc("APPEND NOTES"))
}

func TestUpdatePrependNotesOption(t *testing.T) {
	out, _, code := runThings(t, "", "update", "--auth-token=token", "--id=1", "--prepend-notes=PREPEND NOTES")
	requireSuccess(t, code)
	assertContains(t, out, "prepend-notes="+enc("PREPEND NOTES"))
}

func TestUpdateTagsOption(t *testing.T) {
	out, _, code := runThings(t, "", "update", "--auth-token=token", "--id=1", "--tags=tag1,tag 3,tag2")
	requireSuccess(t, code)
	assertContains(t, out, "tags="+enc("tag1,tag 3,tag2"))
}

func TestUpdateChecklistItemsOption(t *testing.T) {
	out, _, code := runThings(t, "", "update", "--auth-token=token", "--id=1", "--checklist-item=ITEM 1", "--checklist-item=ITEM 2")
	requireSuccess(t, code)
	assertContains(t, out, "checklist-items="+join(enc("ITEM 1"), enc("ITEM 2")))
}

func TestUpdateAppendChecklistItemsOption(t *testing.T) {
	out, _, code := runThings(t, "", "update", "--auth-token=token", "--id=1", "--append-checklist-item=ITEM 1", "--append-checklist-item=ITEM 2")
	requireSuccess(t, code)
	assertContains(t, out, "append-checklist-items="+join(enc("ITEM 1"), enc("ITEM 2")))
}

func TestUpdatePrependChecklistItemsOption(t *testing.T) {
	out, _, code := runThings(t, "", "update", "--auth-token=token", "--id=1", "--prepend-checklist-item=ITEM 1", "--prepend-checklist-item=ITEM 2")
	requireSuccess(t, code)
	assertContains(t, out, "prepend-checklist-items="+join(enc("ITEM 1"), enc("ITEM 2")))
}

func TestUpdateAddTagsOption(t *testing.T) {
	out, _, code := runThings(t, "", "update", "--auth-token=token", "--id=1", "--add-tags=tag1,tag 3,tag2")
	requireSuccess(t, code)
	assertContains(t, out, "add-tags="+enc("tag1,tag 3,tag2"))
}

func TestUpdateAcceptsTitleInput(t *testing.T) {
	out, _, code := runThings(t, "", "update", "--auth-token=token", "--id=1", "New Project")
	requireSuccess(t, code)
	assertContains(t, out, "title="+enc("New Project"))
}

func TestUpdateAcceptsTitleAndNotesInput(t *testing.T) {
	out, _, code := runThings(t, "", "update", "--auth-token=token", "--id=1", "New Todo\n\nThe notes")
	requireSuccess(t, code)
	assertContains(t, out, "title="+enc("New Todo")+"&notes="+enc("The notes"))
}

func TestUpdateAcceptsTitleFromStdin(t *testing.T) {
	out, _, code := runThings(t, "New Todo", "update", "--auth-token=token", "--id=1", "-")
	requireSuccess(t, code)
	assertContains(t, out, "title="+enc("New Todo"))
}

func TestUpdateAcceptsTitleAndNotesFromStdin(t *testing.T) {
	out, _, code := runThings(t, "New Todo\n\nThe notes", "update", "--auth-token=token", "--id=1", "-")
	requireSuccess(t, code)
	assertContains(t, out, "title="+enc("New Todo")+"&notes="+enc("The notes"))
}
