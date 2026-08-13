# things3-cli

[![CI](https://github.com/ossianhempel/things3-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/ossianhempel/things3-cli/actions/workflows/ci.yml)

CLI for Things 3 by Cultured Code, implemented in Go.

This project ships a single Go binary with unit and integration tests.

## Status

Work in progress. The goal is full end-to-end coverage for the Things URL
scheme interactions on macOS.

## Installation (from source)

```
make install
```

## Installation (Homebrew)

```
brew install ossianhempel/tap/things3-cli
```

## Features

- `add`              Add a new todo
- `update`           Update an existing todo (requires auth token)
- `delete`           Delete an existing todo
- `add-area`         Add a new area
- `add-project`      Add a new project
- `update-area`      Update an existing area
- `delete-area`      Delete an existing area
- `update-project`   Update an existing project (requires auth token)
- `list-project-tasks` List tasks for a project ID from the database
- `rename-project`   Rename an existing project (requires auth token)
- `delete-project`   Delete an existing project
- `show`             Show an area, project, tag, or todo from the database (`--recursive` includes checklist items for todos)
- `search`           Search todos and projects in the database
- `inbox`            List inbox tasks
- `today`            List today tasks
- `upcoming`         List upcoming tasks
- `repeating`        List repeating tasks
- `templates`        List repeating template tasks
- `anytime`          List anytime tasks
- `someday`          List someday tasks
- `logbook`          List logbook tasks
- `logtoday`         List tasks completed today
- `createdtoday`     List tasks created today
- `completed`        List completed tasks
- `canceled`         List canceled tasks
- `trash`            List trashed tasks
- `deadlines`        List tasks with deadlines
- `all`              List key sections from the database
- `help`             Command help and man page
- `--version`        Print CLI + Things version info

## Auth token setup (for URL-scheme updates)

Ordinary update operations use the Things URL scheme and require an auth token.
Repeat-only updates write directly to the database and do not require the URL
token. A command that combines ordinary fields with repeat fields requires both
the token and writable database access.

1. Open Things 3.
2. Settings -> General -> Things URLs.
3. Copy the token (or enable "Allow 'things' CLI to access Things").
4. Export it:

```
export THINGS_AUTH_TOKEN=your_token_here
```

Tip: add the export to your shell profile (e.g. `~/.zshrc`) to persist it.
You can run `things auth` to check token status and print these steps.

## Database access

In addition to the URL-scheme commands above, this CLI can read your local
Things database to list content:

- `things projects`  List projects
- `things areas`     List areas
- `things tags`      List tags
- `things tasks`     List todos (with filters)
- `things today`     List Today tasks
- `things templates` List repeating template tasks
- `things list-project-tasks --id <UUID>` List todos for a project

`things today` follows Things' Today/This Evening ordering using the raw
`today_index_reference_date` and `today_index` database fields. Select the
ordering metadata with `--select start_bucket,today_index_reference_date,today_index`.

By default it looks for the Things database in your user Library under the
Things app group container (the `ThingsData-*` folder). You can override the
path with `THINGSDB` or `--db`.

Read commands need database access; repeat commands need writable database
access. Because the database lives inside the Things app sandbox, both normally
require Full Disk Access for your terminal or agent host.

## Repeating todos and projects

Use `--repeat` flags with `add`, `add-project`, `update`, or `update-project`
to create or change repeating templates. These changes write directly to the
Things database, so Full Disk Access is required. Repeating updates require a
single explicit title (for add) or `--id` (for update).

Use `things templates` to list the hidden template rows that control future
instances of recurring todos and projects. Template UUIDs can be passed to
`things update` / `things update-project` when you need to move or edit the
source template rather than the visible generated instance.

Supported patterns are every N day/week/month/year. Use after-completion mode
(the default) when the next copy should be based on completing the current copy;
use schedule mode for a fixed calendar cadence. `--repeat-start` anchors the
recurrence pattern and is separate from the item's ordinary `--when` date.
`--repeat-deadline=N` adds a deadline so each copy appears in Today N days
earlier. Multi-day weekly patterns are not supported yet. Bound a schedule with
either `--repeat-until=YYYY-MM-DD` (stop after a date) or `--repeat-count=N`
(stop after N occurrences); the two are mutually exclusive and count-based
rules omit an end date entirely.

`--dry-run` previews the normalized repeat operation without opening Things or
writing the database. After a real write, the CLI re-reads the template and only
reports verified repeat state. If a multi-stage add partially succeeds, its
output identifies the completed and failed stages and includes a trusted UUID
when available; re-read that UUID before retrying instead of guessing by title.

Note: the CLI verifies the template row after writing. Things itself spawns the
visible current occurrence on its own schedule (typically a nightly pass), so a
verified template may not produce a visible instance until then.

Examples:

```
things add "Daily standup" --repeat=day --repeat-mode=schedule
things update --id <uuid> --repeat=week --repeat-every=2
things update --id <uuid> --repeat-clear
things add-project "Sprint review" --area "Work" --repeat=week --repeat-mode=schedule
things add-project "20 videos in 20 weeks" --repeat=week --repeat-start=2026-08-17 \
  --repeat-count=20 --repeat-deadline=6
things update-project --id <uuid> --repeat=month
things update-project --id <uuid> --repeat-clear
things update --id <uuid> --complete-checklist-item "Book hotel"
things update --id <uuid> --incomplete-checklist-item "Book hotel"
```

## Agent Skills

This repo includes a Things agent skill at `skills/things/SKILL.md`.

That canonical file is mirrored to `../agent-scripts/skills/things/SKILL.md` so
local agent setups can consume the same guidance. Run `make check-things-skill`
to check a sibling checkout and `make sync-things-skill` to update the exact
active mirror safely.

The `skill-sync.yml` workflow runs on every `main` push and compares both
repositories' published `main` artifacts. It uses a read-only deploy key stored
as the `AGENT_SCRIPTS_DEPLOY_KEY` Actions secret to fetch the private
`ossianhempel/agent-scripts` mirror.

Maintainer release guidance lives in `.codex/skills/release-flow/SKILL.md` and
`docs/RELEASING.md`. Use it when preparing a public release or updating the
Homebrew formula.

## Notes

- macOS only (uses the Things URL scheme and `open` under the hood).
- Authentication for update operations follows the Things URL scheme
  authorization model.
- Write commands open Things in the background by default; use `--foreground`
  to bring it to the front, or `--dry-run` to print the URL without opening.
- `show --recursive` includes checklist items in both table and JSON output.
- Delete commands (todo/project/area) use AppleScript and require Things
  automation permission for your terminal (you may see a macOS prompt).
- Delete commands prompt for confirmation when run interactively; pass
  `--confirm` in non-interactive scripts. Use `--dry-run` to preview.
- Aliases: `create-project` -> `add-project`, `create-area` -> `add-area`.
- Scheduling: use `--when=someday` to move to Someday; use `update --later`
  (or `--when=evening`) to move to This Evening.
