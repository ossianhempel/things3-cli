package cli

import (
	"fmt"
	"io"

	"github.com/ossianhempel/things3-cli/internal/db"
)

const previewLimit = 20

func previewTasks(out io.Writer, tasks []db.Task) error {
	count := len(tasks)
	fmt.Fprintf(out, "Matches: %d\n", count)
	if count == 0 {
		return nil
	}
	preview := tasks
	if count > previewLimit {
		preview = tasks[:previewLimit]
		fmt.Fprintf(out, "Preview (first %d):\n", previewLimit)
	} else {
		fmt.Fprintln(out, "Preview:")
	}
	return printTasks(out, preview, TaskOutputOptions{Format: "table"})
}
