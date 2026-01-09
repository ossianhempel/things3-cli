package things

import "strings"

// UpdateProjectOptions defines options for update-project.
type UpdateProjectOptions struct {
	AuthToken      string
	ID             string
	Notes          string
	PrependNotes   string
	AppendNotes    string
	When           string
	Deadline       string
	Tags           string
	AddTags        string
	AreaID         string
	Area           string
	Completed      bool
	Canceled       bool
	Reveal         bool
	Duplicate      bool
	CompletionDate string
	CreationDate   string
	Todos          []string
}

// BuildUpdateProjectURL builds a Things URL for the update-project command.
func BuildUpdateProjectURL(opts UpdateProjectOptions, rawInput string) (string, error) {
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

	params := make([]string, 0, 20)
	params = append(params, "auth-token="+URLEncode(opts.AuthToken))
	params = append(params, "id="+URLEncode(opts.ID))

	if title != "" {
		params = append(params, "title="+URLEncode(title))
	}

	if opts.PrependNotes != "" {
		params = append(params, "prepend-notes="+URLEncode(opts.PrependNotes))
	}

	if opts.AppendNotes != "" {
		params = append(params, "append-notes="+URLEncode(opts.AppendNotes))
	}

	if opts.Duplicate {
		params = append(params, "duplicate=true")
	}

	if opts.When != "" {
		params = append(params, "when="+URLEncode(opts.When))
	}

	if opts.Deadline != "" {
		params = append(params, "deadline="+URLEncode(opts.Deadline))
	}

	if len(opts.Todos) > 0 {
		encoded := make([]string, 0, len(opts.Todos))
		for _, todo := range opts.Todos {
			encoded = append(encoded, URLEncode(todo))
		}
		params = append(params, "to-dos="+Join(encoded...))
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

	if opts.AreaID != "" {
		params = append(params, "area-id="+URLEncode(opts.AreaID))
	} else if opts.Area != "" {
		params = append(params, "area="+URLEncode(opts.Area))
	}

	return "things:///update-project?" + strings.Join(params, "&") + "&", nil
}
