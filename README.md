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
Recurrence requires the native Things Repeat dialog; direct writes to the live database are blocked.

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

Active lists (`today`, `tasks`, `search`, and the other list views) skip to-dos
that are still open inside a trashed, completed, or canceled project, matching
what Things shows. Completing a project archives whatever is left open in it.
Those to-dos remain reachable by UUID with `show --id`, and children that were
closed along with their project still appear in `logbook`, `completed`, and
`canceled`.

By default it looks for the Things database in your user Library under the
Things app group container (the `ThingsData-*` folder). You can override the
path with `THINGSDB` or `--db`.

Read commands need database access; repeat commands need writable database
access. Because the database lives inside the Things app sandbox, both normally
require Full Disk Access for your terminal or agent host.

## Repeating todos and projects

Create and edit recurrence in Things using **Items > Repeat** or **File > New
Repeating To-Do**. For a fixed daily cadence choose daily, every 1 day, and
verify the saved template's repeat settings in the app.

The former SQLite writer could report "verified" merely by reading its own
writes, without Things accepting them. Writes to Things-managed databases are
now rejected before any URL creation/update. This includes explicit `--db`
paths and aliases through parent symlinks. Do not write SQL to the live database.
`things repeating` and `things templates` remain read-only inspection tools;
`--repeat` flags with `--dry-run` remain planning tools, not proof of application
support. Internal write tests use isolated database fixtures.

See [Things' official repeating workflow](https://culturedcode.com/things/support/articles/2803564/).

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
