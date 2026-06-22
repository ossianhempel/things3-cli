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

## Auth token setup (for update commands)

Update operations use the Things URL scheme and require an auth token.

1. Open Things 3.
2. Settings -> General -> Things URLs.
3. Copy the token (or enable "Allow 'things' CLI to access Things").
4. Export it:

```
export THINGS_AUTH_TOKEN=your_token_here
```

Tip: add the export to your shell profile (e.g. `~/.zshrc`) to persist it.
You can run `things auth` to check token status and print these steps.

## Database access (read-only)

In addition to the URL-scheme commands above, this CLI can read your local
Things database to list content:

- `things projects`  List projects
- `things areas`     List areas
- `things tags`      List tags
- `things tasks`     List todos (with filters)
- `things today`     List Today tasks
- `things templates` List repeating template tasks
- `things list-project-tasks --id <UUID>` List todos for a project

By default it looks for the Things database in your user Library under the
Things app group container (the `ThingsData-*` folder). You can override the
path with `THINGSDB` or `--db`.

Note: The database lives inside the Things app sandbox, so you may need to
grant your terminal Full Disk Access.

## Repeating todos

Use `--repeat` flags with `add` or `update`
to create or change repeating templates. These changes write directly to the
Things database, so Full Disk Access is required. Repeating updates require a
single explicit title (for add) or `--id` (for update).

Use `things templates` to list the hidden template rows that control future
instances of recurring todos. Template UUIDs can be passed to `things update`
when you need to move or edit the source template rather than the visible
generated instance.

Supported patterns: every N day/week/month/year, in after-completion (default)
or schedule mode. The anchor date controls weekday/month/day; multi-day weekly
patterns are not supported yet. Use `--repeat-until` to stop after a date.
Repeating projects are not supported.

Examples:

```
things add "Daily standup" --repeat=day --repeat-mode=schedule
things update --id <uuid> --repeat=week --repeat-every=2
things update --id <uuid> --repeat-clear
things update --id <uuid> --complete-checklist-item "Book hotel"
things update --id <uuid> --incomplete-checklist-item "Book hotel"
```

## Scheduling dates

`--when` and `--deadline` accept explicit dates (`YYYY-MM-DD`) and the
existing explicit date-time formats for scheduled reminders. They also accept a
small deterministic set of natural-language dates on `add`, `add-project`,
`update`, and `update-project`: `today`, `tomorrow`, `next <weekday>`,
`next week`, `next month`, `next year`, and `in N day/week/month/year`.
Natural-language dates are resolved in the local timezone and sent to Things as
normalized `YYYY-MM-DD` values.

Examples:

```
things add "Send invoice" --when tomorrow
things add "Renew contract" --deadline "next Friday"
things update --id <uuid> --when "in 2 weeks"
```

Ambiguous phrases are intentionally unsupported. Use an explicit date instead
of phrases like `this Friday`, `later this week`, month names without a year,
or natural-language times. `--when=evening`, `--when=anytime`,
`--when=someday`, and `--when=inbox` remain Things scheduling list values, not
parsed dates.

## Agent Skills

This repo includes a Things agent skill at `skills/things/SKILL.md`.

That file is mirrored to `../agent-scripts/archived-skills/things/SKILL.md` so
local agent setups can consume the same guidance. Keep both copies in sync
when commands or behavior change.

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
- Scheduling: `--when`/`--deadline` accept explicit dates and supported
  natural-language dates such as `tomorrow`, `next Friday`, and `in 2 weeks`.
  Use `--when=someday` to move to Someday; use `update --later` (or
  `--when=evening`) to move to This Evening.
