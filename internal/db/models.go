package db

const (
	StatusIncomplete = 0
	StatusCanceled   = 2
	StatusCompleted  = 3
)

const (
	TaskTypeTodo    = 0
	TaskTypeProject = 1
	TaskTypeHeading = 2
)

type Project struct {
	UUID      string `json:"uuid"`
	Title     string `json:"title"`
	AreaID    string `json:"area_id,omitempty"`
	AreaTitle string `json:"area_title,omitempty"`
	Status    int    `json:"status"`
	Trashed   bool   `json:"trashed"`
}

type Area struct {
	UUID    string `json:"uuid"`
	Title   string `json:"title"`
	Visible bool   `json:"visible"`
}

type Tag struct {
	UUID     string `json:"uuid"`
	Title    string `json:"title"`
	Shortcut string `json:"shortcut,omitempty"`
	ParentID string `json:"parent_id,omitempty"`
	Usage    int    `json:"usage,omitempty"`
}

type Task struct {
	Type                    string          `json:"type,omitempty"`
	UUID                    string          `json:"uuid"`
	Title                   string          `json:"title"`
	Status                  int             `json:"status"`
	Trashed                 bool            `json:"trashed"`
	Notes                   string          `json:"notes,omitempty"`
	Start                   string          `json:"start,omitempty"`
	StartDate               string          `json:"start_date,omitempty"`
	StartBucket             *int            `json:"start_bucket,omitempty"`
	Repeating               bool            `json:"repeating,omitempty"`
	Repeat                  *RepeatState    `json:"repeat,omitempty"`
	Deadline                string          `json:"deadline,omitempty"`
	StopDate                string          `json:"stop_date,omitempty"`
	Created                 string          `json:"created,omitempty"`
	Modified                string          `json:"modified,omitempty"`
	Index                   int             `json:"index,omitempty"`
	TodayIndex              *int            `json:"today_index,omitempty"`
	TodayIndexReferenceDate *int            `json:"today_index_reference_date,omitempty"`
	Tags                    []string        `json:"tags,omitempty"`
	Checklist               []ChecklistItem `json:"checklist,omitempty"`
	ProjectID               string          `json:"project_id,omitempty"`
	ProjectTitle            string          `json:"project_title,omitempty"`
	AreaID                  string          `json:"area_id,omitempty"`
	AreaTitle               string          `json:"area_title,omitempty"`
	HeadingID               string          `json:"heading_id,omitempty"`
	HeadingTitle            string          `json:"heading_title,omitempty"`
}

// RepeatState is the canonical, agent-facing projection of a repeating template.
// Raw recurrence plist data is intentionally excluded.
type RepeatState struct {
	Active         bool   `json:"active"`
	Paused         bool   `json:"paused"`
	Mode           string `json:"mode,omitempty"`
	Unit           string `json:"unit,omitempty"`
	Interval       int    `json:"interval,omitempty"`
	Anchor         string `json:"anchor,omitempty"`
	EndDate        string `json:"end_date,omitempty"`
	DeadlineOffset *int   `json:"deadline_offset,omitempty"`
	Count          *int   `json:"count,omitempty"`
	NextDate       string `json:"next_date,omitempty"`
	Scheduled      bool   `json:"scheduled"`
	StartBucket    *int   `json:"start_bucket,omitempty"`
	DecodeWarning  string `json:"decode_warning,omitempty"`
}

type ChecklistItem struct {
	UUID     string `json:"uuid"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Index    int    `json:"index,omitempty"`
	StopDate string `json:"stop_date,omitempty"`
	Created  string `json:"created,omitempty"`
	Modified string `json:"modified,omitempty"`
}

type Item struct {
	UUID         string          `json:"uuid"`
	Type         string          `json:"type"`
	Title        string          `json:"title"`
	Status       *int            `json:"status,omitempty"`
	Trashed      *bool           `json:"trashed,omitempty"`
	Checklist    []ChecklistItem `json:"checklist,omitempty"`
	ProjectTitle string          `json:"project_title,omitempty"`
	AreaTitle    string          `json:"area_title,omitempty"`
	HeadingTitle string          `json:"heading_title,omitempty"`
	Visible      *bool           `json:"visible,omitempty"`
	Shortcut     string          `json:"shortcut,omitempty"`
	ParentID     string          `json:"parent_id,omitempty"`
}

type TreeItem struct {
	UUID    string     `json:"uuid"`
	Type    string     `json:"type"`
	Title   string     `json:"title"`
	Status  *int       `json:"status,omitempty"`
	Trashed *bool      `json:"trashed,omitempty"`
	Items   []TreeItem `json:"items,omitempty"`
}

type ProjectFilter struct {
	Status         *int
	IncludeTrashed bool
	AreaID         string
}

type TaskFilter struct {
	Status                *int
	IncludeTrashed        bool
	ExcludeTrashedContext bool
	ProjectID             string
	AreaID                string
	TagID                 string
	Search                string
	Limit                 int
	Offset                int
	IncludeChecklist      bool
	Types                 []int
	CreatedBefore         *float64
	CreatedAfter          *float64
	ModifiedBefore        *float64
	ModifiedAfter         *float64
	DueBefore             *int
	StartBefore           *int
	HasURL                *bool
	Order                 string
	IncludeRepeating      bool
	RepeatingOnly         bool
}

func StatusLabel(status int) string {
	switch status {
	case StatusIncomplete:
		return "incomplete"
	case StatusCompleted:
		return "completed"
	case StatusCanceled:
		return "canceled"
	default:
		return "unknown"
	}
}
