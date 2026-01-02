# things3-cli

[![CI](https://github.com/ossianhempel/things3-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/ossianhempel/things3-cli/actions/workflows/ci.yml)

CLI for Things 3 by Cultured Code, implemented in Go.

This project aims to reproduce the behavior and UX of the open source
reference CLI (in another language) while shipping a single Go binary with
unit and integration tests. The reference repository that inspires this work
is:

- https://github.com/itspriddle/things-cli
- https://github.com/thingsapi/things-cli

## Status

Work in progress. The goal is feature parity with the reference CLI and full
end-to-end coverage for the Things URL scheme interactions on macOS.

## Installation (from source)

```
make install
```

## Installation (Homebrew)

```
brew install ossianhempel/tap/things3-cli
```

## Target features (parity with reference)

- `add`              Add a new todo
- `update`           Update an existing todo (requires auth token)
- `add-area`         Add a new area
- `add-project`      Add a new project
- `update-area`      Update an existing area
- `update-project`   Update an existing project (requires auth token)
- `show`             Show an area, project, tag, or todo from the database
- `search`           Search tasks in the database
- `inbox`            List inbox tasks
- `today`            List today tasks
- `upcoming`         List upcoming tasks
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

## Database access (read-only)

In addition to the URL-scheme commands above, this CLI can read your local
Things database to list content:

- `things projects`  List projects
- `things areas`     List areas
- `things tags`      List tags
- `things tasks`     List todos (with filters)
- `things today`     List Today tasks

By default it looks for the Things database in your user Library under the
Things app group container (the `ThingsData-*` folder). You can override the
path with `THINGSDB` or `--db`.

Note: The database lives inside the Things app sandbox, so you may need to
grant your terminal Full Disk Access.

## Notes

- macOS only (uses the Things URL scheme and `open` under the hood).
- Authentication for update operations follows the Things URL scheme
  authorization model.
- Write commands open Things in the background by default; use `--foreground`
  to bring it to the front, or `--dry-run` to print the URL without opening.
- For update operations you can set `THINGS_AUTH_TOKEN` to avoid passing
  `--auth-token` every time.
- Area creation/updates use AppleScript and require Things automation
  permission for your terminal (you may see a macOS prompt).
- Aliases: `create-project` -> `add-project`, `create-area` -> `add-area`.
- Scheduling: use `--when=someday` to move to Someday; use `update --later`
  (or `--when=evening`) to move to This Evening.
