# Repository Guidelines

## Project Structure & Module Organization
- `cmd/things/` is the CLI entry point.
- `internal/cli/` defines commands, help text, and output formatting.
- `internal/db/` contains Things database models, queries, and read-only access.
- `internal/things/` implements URL scheme encoding and add/update logic.
- `internal/osascript/` wraps AppleScript helpers for area/project automation.
- `integration/` holds integration tests and fixtures; unit tests live alongside code in `internal/**`.
- `docs/`, `doc/`, and `share/` include documentation and man page sources; `Formula/` contains the Homebrew formula.
- `scripts/` contains release automation helpers.

## Build, Test, and Development Commands
- `make build` — builds the CLI binary into `bin/things`.
- `make test` — runs all unit and integration tests (`go test ./...`).
- `make install` — builds and installs the binary (and man page if present) to `/usr/local` or `PREFIX`.
- `make uninstall` — removes installed binary and man page.
- Example direct Go commands: `go test ./integration`, `go build -o bin/things ./cmd/things`.

## Coding Style & Naming Conventions
- Go code follows standard `gofmt` formatting (tabs for indentation).
- Package names are lower-case, short, and domain-focused (`cli`, `db`, `things`).
- Command files in `internal/cli/` mirror subcommand names (e.g., `add.go`, `update.go`).
- Tests use the Go convention `*_test.go`.

## Testing Guidelines
- Tests use Go’s built-in `testing` package.
- Unit tests live next to their packages (e.g., `internal/db/*_test.go`).
- Integration tests live in `integration/` and include sqlite fixtures and helpers.
- Run the full suite with `make test` or target a package with `go test ./internal/cli -run TestName`.
- When Things 3 is available locally, verify user-facing CLI changes against the
  real app/database in addition to fixtures. Prefer privacy-light smoke tests
  for read-only commands, such as selecting only UUIDs or limiting output.
- For write-path changes, create clearly named temporary Things items (for
  example, prefixed with `things3-cli test`) that exercise the behavior, report
  exactly what was created and verified, and then delete, complete, or otherwise
  clean them up unless the developer explicitly wants to inspect them first.
- If live verification is skipped because Things, Full Disk Access, automation
  permission, or `THINGS_AUTH_TOKEN` is unavailable, call that out in the final
  testing notes.

## Commit & Pull Request Guidelines
- Recent commits use short, imperative summaries; common prefixes include `feat:`, `chore:`, and `docs:`.
- If you use a prefix, keep it concise and match the scope of the change.
- No PR template is defined; include a brief summary and testing notes.
- Call out macOS/Things-specific behavior changes and any permissions needed.

## Configuration & Permissions Notes
- Update commands require `THINGS_AUTH_TOKEN`; see README for setup.
- DB reads may require Full Disk Access for your terminal.
- Area/project automation uses AppleScript and may trigger macOS prompts.

## Agent Scripts / Skills Sync
- Treat `skills/things/SKILL.md` as canonical and keep it in sync with the
  active mirror at `../agent-scripts/skills/things/SKILL.md`.
- After CLI changes, update both copies so agent guidance stays aligned.
- Use `make check-things-skill` to check parity and `make sync-things-skill`
  for the symlink-safe, atomic local sync path.
- Use `.codex/skills/release-flow/SKILL.md` for public release prep,
  validation, Homebrew formula updates, and GitHub release verification.
- If release mechanics change, update `.codex/skills/release-flow/SKILL.md`,
  `.codex/skills/release-flow/references/release-flow.md`, and
  `docs/RELEASING.md` together.
