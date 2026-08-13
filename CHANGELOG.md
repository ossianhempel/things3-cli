# Changelog

All notable changes to this project will be documented in this file.

The format is based on Keep a Changelog, and this project adheres to
Semantic Versioning.

## [Unreleased]
- Added `--repeat-count=N` to stop a repeating schedule after N occurrences,
  mutually exclusive with `--repeat-until`. Count-based rules omit an end date
  entirely, matching Things' own encoding (`rc` in the recurrence rule).
- Fixed repeat anchors with a deadline: `ia` is now shifted by the deadline
  offset so occurrences land on the anchor weekday, and a future anchor is
  itself the first instance.
- Added repeating project support: `add-project --repeat=...` and
  `update-project --repeat=...` create or clear repeating project templates,
  matching the repeat semantics already available for todos. Repeat flags on
  projects support after-completion or schedule mode, every N units, an anchor
  date, a stop date, a repeat count, and repeating deadlines.
- `things templates` and `things repeating` now include repeating projects, not
  just repeating todos.
- Documented that Things spawns the visible current occurrence on its own
  schedule (typically a nightly pass) after a repeat template is verified.

## [0.4.0] - 2026-07-18
- Fixed repeating todos remaining paused as Someday templates after applying a repeat rule.
- Added verified, stage-aware results for repeat add, update, clear, and dry-run operations.
- Added structured JSON output, shell-safe recovery commands, and database-target provenance for repeat operations.
- Exposed canonical repeat semantics and activation state through the `templates` command.
- Clarified repeat modes, permissions, unsupported project recurrence, and the UUID-first agent workflow.
- Added safe cross-repository synchronization and drift checks for the Things skill.

## [0.3.0] - 2026-05-26
- Added `templates` command to list repeating template tasks from the Things database.
- Added `update --complete-checklist-item` and `--incomplete-checklist-item`
  to change existing checklist item completion status via the Things JSON URL
  endpoint.
- Clear task and project notes by passing an empty `--notes` value to `update` and `update-project`.
- Include projects in `search` results.

## [0.2.1] - 2026-04-20
- Accept `ThingsData-*` and `Things Database.thingsdatabase` directories in `--db` and `THINGSDB`.
- Return a clear error when a directory does not contain a Things database file.
- Added checklist output support for `show --recursive`.
- Added `update-project`, `list-project-tasks`, and `rename-project` commands.
- Exposed `start_bucket` in Today task output.
- Hardened GitHub Actions and AppleScript string interpolation.
- Synced README, root help, and in-repo agent skill guidance with current commands.

## [0.2.0] - 2026-01-09
- Added guardrails for unsafe titles (e.g. tag=work) with --allow-unsafe-title override.
- Require auth token before URL updates; error early with clearer messaging.
- Verify --when/--later updates against the database to avoid false positives (opt-out with --no-verify).
- Prevent moving non-today tasks to This Evening unless --allow-non-today is set.
- Require confirmation for query deletes (prompt or --confirm=delete/--yes).

## [0.1.0] - 2026-01-06
- Initial Go port of `things-cli` (commands, help, man page, tests).
- Added read-only database commands (`projects`, `areas`, `tags`, `tasks`).
- Fix repeating add to preserve scheduling fields so templates are not trashed.
