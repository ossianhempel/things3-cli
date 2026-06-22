---
name: things
description: Things 3 CLI for reading the local Things database and creating/updating tasks/projects via the Things URL scheme. Use when you need to list, search, or modify Things tasks, projects, areas, or tags from the terminal (deletes use AppleScript).
---

# things

Use `things` to read from the local Things 3 database and to create/update items via the Things URL scheme (areas use AppleScript).

Quick start (read)
- `things inbox`
- `things today`
- `things repeating`
- `things templates`
- `things projects` / `things areas` / `things tags`
- `things tasks --search "query" --json`
- `things tasks --query 'tag:work AND title:/review/i' --format jsonl`
- `things search --query 'tag:PLAN' --select type,title` searches both todos and projects.
- `things show --project "Project Name"`
- `things show --id <TODO_UUID> --recursive --json`

Write (URL scheme)
- `things add "Task title" --notes "..." --list "Project or Area"`
- Checklist items (repeat `--checklist-item` per item): `things add "My task" --checklist-item="Step 1" --checklist-item="Step 2" --checklist-item="Step 3"`
- `things add-project "Project title" --area "Area Name"`
- `things update --id <uuid> --notes "Updated notes"`
- Clear notes by passing an explicit empty notes value: `things update --id <uuid> --notes ""`
- Clear a deadline, scheduled date, or all tags the same way: `things update --id <uuid> --deadline=""`, `--when=""`, or `--tags=""` (also works for `update-project`)
- Natural scheduling dates: `--when`/`--deadline` accept explicit dates plus `today`, `tomorrow`, `next <weekday>`, `next week/month/year`, and `in N day/week/month/year`; dry-run shows the normalized `YYYY-MM-DD` value.
- Checklist status by exact title: `things update --id <uuid> --complete-checklist-item "Step 1"` or `things update --id <uuid> --incomplete-checklist-item "Step 1"`
- Bulk update (preview then apply): `things update --query 'tag:work' --dry-run` then `things update --query 'tag:work' --yes --tags "Work"`
- `things update-project --id <uuid> "New project title"`
- `things rename-project --id <uuid> --title "New project title"`
- `things list-project-tasks --id <project_uuid>`
- `things delete --id <uuid>` or `things delete "Todo title"`
- Bulk delete (preview then apply): `things delete --query 'notes:/deprecated/i' --dry-run` then `things delete --query 'notes:/deprecated/i' --yes`
- Undo last bulk action: `things undo --dry-run` then `things undo --yes`
- `things delete-project --id <uuid>` or `things delete-project "Project title"`
- `things delete-area --id <uuid>` or `things delete-area "Area Name"`
- Move to Someday: `things update --id <uuid> --when=someday`
- Move to This Evening (Later): `things update --id <uuid> --later` (alias for `--when=evening`)

Repeating (database writes)
- Create repeating todo: `things add "Daily standup" --repeat=day --repeat-mode=schedule`
- Repeat every 2 weeks after completion: `things update --id <uuid> --repeat=week --repeat-every=2`
- Anchor a schedule on a date: `things add "Monthly bill" --repeat=month --repeat-mode=schedule --repeat-start=2026-02-01`
- Add repeating deadlines: `things update --id <uuid> --repeat=week --repeat-deadline=2`
- Stop repeating after a date: `things update --id <uuid> --repeat=day --repeat-until=2027-06-01`
- Clear repeating schedule: `things update --id <uuid> --repeat-clear`
 - Repeating adds launch Things in the background first to ensure the item hits the database before the repeat rule is applied.

Filters + DB
- Use `--db` or `THINGSDB` to point to a specific Things database. Accepted forms: the `main.sqlite` file, the `Things Database.thingsdatabase` directory, or the parent `ThingsData-*` directory.
- Common filters: `--filter-project`, `--filter-area`, `--filter-tag`, `--status`, `--search`, `--limit`, `--offset`.
- Rich query: `--query` supports boolean ops, field predicates, and regex (e.g. `title:/regex/ AND tag:work`).
- Repeating tasks: `things repeating` or `--query 'repeating:true'`.
- Repeating templates: `things templates --area "Area Name"` lists hidden template rows that control future recurring instances.
- Date filters: `--created-before/after`, `--modified-before/after`, `--due-before`, `--start-before`.
- URL filter: `--has-url`.
- Sorting: `--sort created,-deadline,title`.
- Output: `--format table|json|jsonl|csv`, `--select uuid,title,status`, `--no-header`. `--json` still works.
- `--recursive` includes checklist items in JSON output.

Auth + permissions
- Updates require an auth token: run `things auth` for setup/status, set `THINGS_AUTH_TOKEN`, or pass `--auth-token`.
- DB reads may require Full Disk Access for your terminal.
- Repeating add/update writes directly to the Things database and requires Full Disk Access; use `--db` or `THINGSDB` to target a specific database.
- URL scheme writes can open/foreground Things; use `--dry-run` to print URLs or `--foreground` to force focus.
- Update `--when/--later` is verified against the database by default; use `--no-verify` to skip verification.
- `--later` / `--when=evening` refuses to move tasks that are already scheduled for a non-today date; use `--allow-non-today` to override.
- Titles that look like flag assignments (e.g. `tag=work`) are rejected; pass `--allow-unsafe-title` to keep them as the title.
- For updates, an explicitly provided empty value clears the field (`--notes ""`, `--deadline=""`, `--when=""`, `--tags=""`); omitting the flag leaves the field unchanged.
- Delete commands prompt for confirmation when interactive; pass `--confirm` for single deletes in non-interactive scripts. Query deletes require `--confirm=delete` or `--yes`.
- Bulk update/delete write an action log; use `things undo` to revert the last bulk update or trash.

Notes
- macOS only.
- Use `things --help` and `things <command> --help` for the full flag list.
- Install via Homebrew: `brew install ossianhempel/tap/things3-cli`.
- When working inside the repo, prefer the repo binary (`make build` then `./bin/things`) or `go run ./cmd/things` to pick up local changes.
