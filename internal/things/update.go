package things

import "strings"

// UpdateOptions defines options for update.
type UpdateOptions struct {
	AuthToken             string
	ID                    string
	Notes                 string
	PrependNotes          string
	AppendNotes           string
	When                  string
	Later                 bool
	Deadline              string
	Tags                  string
	AddTags               string
	Completed             bool
	Canceled              bool
	Reveal                bool
	Duplicate             bool
	CompletionDate        string
	CreationDate          string
	Heading               string
	List                  string
	ListID                string
	ChecklistItems        []string
	PrependChecklistItems []string
	AppendChecklistItems  []string
}

// BuildUpdateURL builds a Things URL for the update command.
func BuildUpdateURL(opts UpdateOptions, rawInput string) (string, error) {
	if opts.AuthToken == "" {
		return "", ErrMissingAuthToken
	}
	if opts.ID == "" {
		return "", errMissingID
	}

	var title string
	notes := opts.Notes

	if rawInput != "" {
		if HasMultipleLines(rawInput) {
			title = FindTitle(rawInput)
			notes = FindNotes(rawInput)
		} else {
			title = rawInput
		}
	}

	params := make([]string, 0, 24)
	params = append(params, "auth-token="+URLEncode(opts.AuthToken))
	params = append(params, "id="+URLEncode(opts.ID))

	if title != "" {
		params = append(params, "title="+URLEncode(title))
	}

	if len(opts.ChecklistItems) > 0 {
		encoded := make([]string, 0, len(opts.ChecklistItems))
		for _, item := range opts.ChecklistItems {
			encoded = append(encoded, URLEncode(item))
		}
		params = append(params, "checklist-items="+Join(encoded...))
	}

	if len(opts.PrependChecklistItems) > 0 {
		encoded := make([]string, 0, len(opts.PrependChecklistItems))
		for _, item := range opts.PrependChecklistItems {
			encoded = append(encoded, URLEncode(item))
		}
		params = append(params, "prepend-checklist-items="+Join(encoded...))
	}

	if len(opts.AppendChecklistItems) > 0 {
		encoded := make([]string, 0, len(opts.AppendChecklistItems))
		for _, item := range opts.AppendChecklistItems {
			encoded = append(encoded, URLEncode(item))
		}
		params = append(params, "append-checklist-items="+Join(encoded...))
	}

	if opts.PrependNotes != "" {
		params = append(params, "prepend-notes="+URLEncode(opts.PrependNotes))
	}

	if opts.AppendNotes != "" {
		params = append(params, "append-notes="+URLEncode(opts.AppendNotes))
	}

	if opts.Heading != "" {
		params = append(params, "heading="+URLEncode(opts.Heading))
	}

	if opts.Duplicate {
		params = append(params, "duplicate=true")
	}

	if opts.When != "" {
		params = append(params, "when="+URLEncode(opts.When))
	} else if opts.Later {
		params = append(params, "when=evening")
	}

	if opts.Deadline != "" {
		params = append(params, "deadline="+URLEncode(opts.Deadline))
	}

	if opts.Reveal {
		params = append(params, "reveal=true")
	}

	if opts.Tags != "" {
		params = append(params, "tags="+URLEncode(opts.Tags))
	}

	if opts.AddTags != "" {
		params = append(params, "add-tags="+URLEncode(opts.AddTags))
	}

	if notes != "" {
		params = append(params, "notes="+URLEncode(notes))
	}

	if opts.CreationDate != "" {
		params = append(params, "creation-date="+URLEncode(opts.CreationDate))
	}

	if opts.CompletionDate != "" {
		params = append(params, "completion-date="+URLEncode(opts.CompletionDate))
	}

	if opts.Canceled {
		params = append(params, "canceled=true")
	} else if opts.Completed {
		params = append(params, "completed=true")
	}

	if opts.ListID != "" {
		params = append(params, "list-id="+URLEncode(opts.ListID))
	} else if opts.List != "" {
		params = append(params, "list="+URLEncode(opts.List))
	}

	return "things:///update?" + strings.Join(params, "&") + "&", nil
}
