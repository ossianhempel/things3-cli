package integration_test

import "testing"

func TestUpdateProjectAuthTokenRequired(t *testing.T) {
	_, errOut, code := runThings(t, "", "update-project")
	requireFailure(t, code)
	assertContains(t, errOut, "Missing Things auth token")
}

func TestUpdateProjectIDRequired(t *testing.T) {
	_, errOut, code := runThings(t, "", "update-project", "--auth-token=token")
	requireFailure(t, code)
	assertContains(t, errOut, "Must specify --id=id")
}

func TestUpdateProjectAreaIDOption(t *testing.T) {
	out, _, code := runThings(t, "", "update-project", "--auth-token=token", "--id=1", "--area-id=1")
	requireSuccess(t, code)
	assertContains(t, out, "area-id="+enc("1"))
}

func TestUpdateProjectAreaOption(t *testing.T) {
	out, _, code := runThings(t, "", "update-project", "--auth-token=token", "--id=1", "--area=SOME AREA")
	requireSuccess(t, code)
	assertContains(t, out, "area="+enc("SOME AREA"))
}

func TestUpdateProjectAreaIDPrecedence(t *testing.T) {
	out, _, code := runThings(t, "", "update-project", "--auth-token=token", "--id=1", "--area-id=1", "--area=SOME AREA")
	requireSuccess(t, code)
	assertContains(t, out, "area-id="+enc("1"))
	assertNotContains(t, out, "area="+enc("SOME AREA"))
}

func TestUpdateProjectCanceledOptions(t *testing.T) {
	out, _, code := runThings(t, "", "update-project", "--auth-token=token", "--id=1", "--canceled")
	requireSuccess(t, code)
	assertContains(t, out, "canceled=true")

	out, _, code = runThings(t, "", "update-project", "--auth-token=token", "--id=1", "--cancelled")
	requireSuccess(t, code)
	assertContains(t, out, "canceled=true")
}

func TestUpdateProjectCanceledPrecedence(t *testing.T) {
	out, _, code := runThings(t, "", "update-project", "--auth-token=token", "--id=1", "--canceled", "--completed")
	requireSuccess(t, code)
	assertContains(t, out, "canceled=true")
	assertNotContains(t, out, "completed=true")
}

func TestUpdateProjectCompletedOption(t *testing.T) {
	out, _, code := runThings(t, "", "update-project", "--auth-token=token", "--id=1", "--completed")
	requireSuccess(t, code)
	assertContains(t, out, "completed=true")

	out, _, code = runThings(t, "", "update-project", "--auth-token=token", "--id=1", "--cancelled")
	requireSuccess(t, code)
	assertContains(t, out, "canceled=true")
}

func TestUpdateProjectCompletionDateOption(t *testing.T) {
	out, _, code := runThings(t, "", "update-project", "--auth-token=token", "--id=1", "--completion-date=2021-05-20")
	requireSuccess(t, code)
	assertContains(t, out, "completion-date=2021-05-20")
}

func TestUpdateProjectDeadlineOption(t *testing.T) {
	out, _, code := runThings(t, "", "update-project", "--auth-token=token", "--id=1", "--deadline=2021-05-20")
	requireSuccess(t, code)
	assertContains(t, out, "deadline=2021-05-20")
}

func TestUpdateProjectNotesOption(t *testing.T) {
	out, _, code := runThings(t, "", "update-project", "--auth-token=token", "--id=1", "--notes=NOTES NOTES")
	requireSuccess(t, code)
	assertContains(t, out, "notes="+enc("NOTES NOTES"))
}

func TestUpdateProjectAppendNotesOption(t *testing.T) {
	out, _, code := runThings(t, "", "update-project", "--auth-token=token", "--id=1", "--append-notes=APPEND NOTES")
	requireSuccess(t, code)
	assertContains(t, out, "append-notes="+enc("APPEND NOTES"))
}

func TestUpdateProjectPrependNotesOption(t *testing.T) {
	out, _, code := runThings(t, "", "update-project", "--auth-token=token", "--id=1", "--prepend-notes=PREPEND NOTES")
	requireSuccess(t, code)
	assertContains(t, out, "prepend-notes="+enc("PREPEND NOTES"))
}

func TestUpdateProjectRevealOption(t *testing.T) {
	out, _, code := runThings(t, "", "update-project", "--auth-token=token", "--id=1", "--reveal")
	requireSuccess(t, code)
	assertContains(t, out, "reveal=true")
}

func TestUpdateProjectTagsOption(t *testing.T) {
	out, _, code := runThings(t, "", "update-project", "--auth-token=token", "--id=1", "--tags=tag1,tag 3,tag2")
	requireSuccess(t, code)
	assertContains(t, out, "tags="+enc("tag1,tag 3,tag2"))
}

func TestUpdateProjectAddTagsOption(t *testing.T) {
	out, _, code := runThings(t, "", "update-project", "--auth-token=token", "--id=1", "--add-tags=tag1,tag 3,tag2")
	requireSuccess(t, code)
	assertContains(t, out, "add-tags="+enc("tag1,tag 3,tag2"))
}

func TestUpdateProjectWhenOption(t *testing.T) {
	out, _, code := runThings(t, "", "update-project", "--auth-token=token", "--id=1", "--when=2021-05-20")
	requireSuccess(t, code)
	assertContains(t, out, "when=2021-05-20")
}

func TestUpdateProjectTodoOption(t *testing.T) {
	out, _, code := runThings(t, "", "update-project", "--auth-token=token", "--id=1", "--todo=Todo 1", "--todo=Todo 2")
	requireSuccess(t, code)
	assertContains(t, out, "to-dos="+join(enc("Todo 1"), enc("Todo 2")))
}

func TestUpdateProjectAcceptsTitleInput(t *testing.T) {
	out, _, code := runThings(t, "", "update-project", "--auth-token=token", "--id=1", "New Project")
	requireSuccess(t, code)
	assertContains(t, out, "title="+enc("New Project"))
}

func TestUpdateProjectAcceptsTitleAndNotesInput(t *testing.T) {
	out, _, code := runThings(t, "", "update-project", "--auth-token=token", "--id=1", "New Project\n\nThe notes")
	requireSuccess(t, code)
	assertContains(t, out, "title="+enc("New Project")+"&notes="+enc("The notes"))
}

func TestUpdateProjectAcceptsTitleFromStdin(t *testing.T) {
	out, _, code := runThings(t, "New Project", "update-project", "--auth-token=token", "--id=1", "-")
	requireSuccess(t, code)
	assertContains(t, out, "title="+enc("New Project"))
}

func TestUpdateProjectAcceptsTitleAndNotesFromStdin(t *testing.T) {
	out, _, code := runThings(t, "New Project\n\nThe notes", "update-project", "--auth-token=token", "--id=1", "-")
	requireSuccess(t, code)
	assertContains(t, out, "title="+enc("New Project")+"&notes="+enc("The notes"))
}
