package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/ossianhempel/things3-cli/internal/db"
)

type TaskSection struct {
	Title string    `json:"title"`
	Items []db.Task `json:"items"`
}

func printProjects(out io.Writer, projects []db.Project, asJSON bool, noHeader bool) error {
	if asJSON {
		enc := json.NewEncoder(out)
		return enc.Encode(projects)
	}
	w := tabwriter.NewWriter(out, 0, 2, 2, ' ', 0)
	if !noHeader {
		fmt.Fprintln(w, "UUID\tTITLE\tAREA\tSTATUS\tTRASHED")
	}
	for _, p := range projects {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%t\n", p.UUID, p.Title, p.AreaTitle, db.StatusLabel(p.Status), p.Trashed)
	}
	return w.Flush()
}

func printAreas(out io.Writer, areas []db.Area, asJSON bool, noHeader bool) error {
	if asJSON {
		enc := json.NewEncoder(out)
		return enc.Encode(areas)
	}
	w := tabwriter.NewWriter(out, 0, 2, 2, ' ', 0)
	if !noHeader {
		fmt.Fprintln(w, "UUID\tTITLE\tVISIBLE")
	}
	for _, a := range areas {
		fmt.Fprintf(w, "%s\t%s\t%t\n", a.UUID, a.Title, a.Visible)
	}
	return w.Flush()
}

func printTags(out io.Writer, tags []db.Tag, asJSON bool, noHeader bool) error {
	if asJSON {
		enc := json.NewEncoder(out)
		return enc.Encode(tags)
	}
	w := tabwriter.NewWriter(out, 0, 2, 2, ' ', 0)
	if !noHeader {
		fmt.Fprintln(w, "UUID\tTITLE\tSHORTCUT\tPARENT")
	}
	for _, tag := range tags {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", tag.UUID, tag.Title, tag.Shortcut, tag.ParentID)
	}
	return w.Flush()
}

func printTasks(out io.Writer, tasks []db.Task, opts TaskOutputOptions) error {
	return writeTasks(out, tasks, opts)
}

func printTaskSections(out io.Writer, sections []TaskSection, opts TaskOutputOptions) error {
	if opts.Format == "json" {
		enc := json.NewEncoder(out)
		return enc.Encode(sections)
	}
	for i, section := range sections {
		if i > 0 {
			fmt.Fprintln(out)
		}
		fmt.Fprintln(out, section.Title)
		if len(section.Items) == 0 {
			continue
		}
		if err := printTasks(out, section.Items, TaskOutputOptions{Format: "table", NoHeader: opts.NoHeader, Select: opts.Select}); err != nil {
			return err
		}
	}
	return nil
}

func printItem(out io.Writer, item *db.Item, asJSON bool, noHeader bool) error {
	if asJSON {
		enc := json.NewEncoder(out)
		return enc.Encode(item)
	}
	w := tabwriter.NewWriter(out, 0, 2, 2, ' ', 0)
	if !noHeader {
		fmt.Fprintln(w, "TYPE\tUUID\tTITLE\tSTATUS\tTRASHED\tPROJECT\tAREA\tHEADING\tVISIBLE\tSHORTCUT\tPARENT")
	}
	status := ""
	if item.Status != nil {
		status = db.StatusLabel(*item.Status)
	}
	trashed := ""
	if item.Trashed != nil {
		trashed = fmt.Sprintf("%t", *item.Trashed)
	}
	visible := ""
	if item.Visible != nil {
		visible = fmt.Sprintf("%t", *item.Visible)
	}
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
		item.Type, item.UUID, item.Title, status, trashed, item.ProjectTitle, item.AreaTitle, item.HeadingTitle, visible, item.Shortcut, item.ParentID)
	return w.Flush()
}

func printTree(out io.Writer, items []db.TreeItem, asJSON bool) error {
	if asJSON {
		enc := json.NewEncoder(out)
		return enc.Encode(items)
	}
	printTreeItems(out, items, "")
	return nil
}

func printTreeItems(out io.Writer, items []db.TreeItem, indent string) {
	for _, item := range items {
		line := fmt.Sprintf("%s- %s", indent, item.Title)
		if item.Type != "" {
			line = fmt.Sprintf("%s (%s)", line, item.Type)
		}
		fmt.Fprintln(out, line)
		if len(item.Items) > 0 {
			printTreeItems(out, item.Items, indent+"  ")
		}
	}
}
