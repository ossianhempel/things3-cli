---
name: things
description: Things 3 CLI for reading the local Things database and creating/updating tasks/projects via the Things URL scheme. Use when you need to list, search, or modify Things tasks, projects, areas, or tags from the terminal (deletes use AppleScript).
---

# things

Use `things` to read from the local Things 3 database and to create/update items via the Things URL scheme (areas use AppleScript).

## Safe operating workflow

For any write whose target is not already a trusted UUID:

1. **Read** with `things search`, `things tasks`, or `things templates`.
2. **Identify** the intended row and capture its UUID. If multiple items match, show the candidates and ask the human; never guess by title.
3. **Preview** the exact write with `--dry-run` when the command supports it.
4. **Write** using `--id <UUID>`.
5. **Verify** by re-reading that UUID; recurrence additionally requires inspecting the saved native Repeat settings.
6. **Report** what was requested, what was verified, and any remaining uncertainty.

Reads do not authorize writes. Prefer UUIDs even when a command accepts a title.

Quick start (read)
- `things inbox`
- `things today` follows the app's Today/This Evening order; select `start_bucket,today_index_reference_date,today_index` to inspect its raw ordering metadata.
- `things repeating`
- `things templates`
- `things projects` / `things areas` / `things tags`
- `things tasks --search "query" --json`
- `things tasks --query 'tag:work AND title:/review/i' --format jsonl`
- `things search --query 'tag:PLAN' --select type,title` searches both todos and projects.
- `things show --project "Project Name"`
- `things show --id <TODO_UUID> --recursive --json`

Completed / done tasks (active lists like `tasks`, `today`, `search` hide closed items by default)
- `things logtoday` lists to-dos completed or canceled today.
- `things logbook` lists completed and canceled to-dos, newest first; `things completed` / `things canceled` narrow to one status.
- Filter by completion date: `things logbook --completed-after 2026-09-28 --completed-before 2026-10-03 --select title,status,stop_date,project`. Date-only `--completed-before` is inclusive of that day.
- Sort by completion date with `--sort -completed` (newest first) or `--sort completed`; `stop_date` is the selectable completion timestamp.
- To include closed items in a general query, use `things tasks --status completed` (or `canceled`/`any`) or `--all`.

Write (URL scheme)
- `things add "Task title" --notes "..." --list "Project or Area"`
- Checklist items (repeat `--checklist-item` per item): `things add "My task" --checklist-item="Step 1" --checklist-item="Step 2" --checklist-item="Step 3"`
- `things add-project "Project title" --area "Area Name"`
- `things update --id <uuid> --notes "Updated notes"`
- Clear notes by passing an explicit empty notes value: `things update --id <uuid> --notes ""`
- Clear a deadline, scheduled date, or all tags the same way: `things update --id <uuid> --deadline=""`, `--when=""`, or `--tags=""` (also works for `update-project`)
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

Repeating tasks and projects (native Things UI)

Use Things itself to create or change recurrence. The URL scheme does not expose repeat creation. The CLI's former direct SQLite writer could report success while Things still treated the item as a normal task; database read-back is not application verification. Never mutate the live database with SQL, including when Things is closed. Never explain a missing repeat rule as a nightly-generation delay.

1. Read existing items and their parent area/project; retain UUIDs to avoid duplicates. For new items, create an ordinary task in the intended parent via the URL scheme, then locate it.
2. Open the identified item in Things using Quick Find or its Things link. Use native Items > Repeat. For a new template, File > New Repeating To-Do is also available.
3. For "once a day", choose **daily**, every **1** day, starting today, ending **never**. Choose **after completion** only when that is the requested behavior.
4. Save with OK. Reopen the template's repeat settings and verify the cadence, start date, end condition, title, and parent in Things. A generated occurrence and its template are different items; retain both UUIDs when available.
5. Use `things templates` / `things repeating` as a supplemental read. Report success only after the native saved rule is visible. If UI access is unavailable, explain the concrete blocker and leave recurrence unconfirmed.

The CLI refuses writes to Things-managed database paths. Repeat flags remain useful for `--dry-run` planning and isolated database fixtures; they are not a supported live recurrence backend. Do not retry with `--db`, move the database, or substitute a raw SQLite script to bypass this limit.

Official workflow: https://culturedcode.com/things/support/articles/2803564/

Filters + DB
- Use `--db` or `THINGSDB` to point to a specific Things database. Accepted forms: the `main.sqlite` file, the `Things Database.thingsdatabase` directory, or the parent `ThingsData-*` directory.
- Common filters: `--filter-project`, `--filter-area`, `--filter-tag`, `--status`, `--search`, `--limit`, `--offset`.
- Rich query: `--query` supports boolean ops, field predicates, and regex (e.g. `title:/regex/ AND tag:work`).
- Repeating tasks and projects: `things repeating` or `--query 'repeating:true'`.
- Repeating templates: `things templates --area "Area Name"` lists hidden template rows that control future recurring instances.
- Date filters: `--created-before/after`, `--modified-before/after`, `--completed-before/after`, `--due-before`, `--start-before`.
- URL filter: `--has-url`.
- Sorting: `--sort created,-deadline,title`; `completed` sorts by completion date.
- Output: `--format table|json|jsonl|csv`, `--select uuid,title,status`, `--no-header`. `--json` still works. `today_index_reference_date` is selectable as a raw packed integer but is not a generic `--sort` field.
- `--recursive` includes checklist items in JSON output.
- Active lists skip to-dos left open inside a trashed, completed, or canceled project, matching Things. Use `show --id <UUID>` to reach one anyway, or `logbook`/`completed`/`canceled` for children closed with their project.

Auth + permissions
- Read-only database commands may require Full Disk Access for the terminal or agent host.
- Ordinary URL-scheme updates require an auth token: run `things auth`, set `THINGS_AUTH_TOKEN`, or pass `--auth-token`.
- Recurrence requires the native Things UI; Full Disk Access only enables reads and does not make direct SQLite writes a supported app interface.
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
