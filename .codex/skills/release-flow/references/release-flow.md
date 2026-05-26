# things3-cli Release Flow

This repo is consumed by external users and agents. Treat releases as public API changes.

## Preconditions

- Work from `main`.
- Check `git status --short`; release scripts require a clean tree for publish.
- Confirm the next version from `CHANGELOG.md`, `git tag --sort=-version:refname`, and the scope of the change.
- Prefer SemVer:
  - patch for bug fixes and docs-only release support,
  - minor for new flags, commands, or output fields,
  - major for breaking command/output behavior.

## Required Validation

Run these before publishing:

```sh
make test
./scripts/build-release.sh vX.Y.Z
./scripts/release-notes.sh vX.Y.Z
```

Check the generated artifacts:

```sh
ls -la dist
cat dist/checksums.txt
tmpdir=$(mktemp -d)
tar -xzf dist/things-X.Y.Z-darwin-arm64.tar.gz -C "$tmpdir"
"$tmpdir/things" --version
rm -rf "$tmpdir"
```

## Changelog

Move `CHANGELOG.md` entries out of `Unreleased` into:

```md
## [X.Y.Z] - YYYY-MM-DD
```

The GitHub release body must be only the bullets for that version. `./scripts/release-notes.sh vX.Y.Z` is the source of truth.

## Homebrew Formula

Update the formula after artifacts exist:

```sh
./scripts/update-brew-formula.sh --version vX.Y.Z --tap-dir ~/Developer/homebrew-tap
```

Validate:

```sh
rg -n 'version "|download/v|sha256' Formula/things3-cli.rb ~/Developer/homebrew-tap/Formula/things3-cli.rb
```

After the GitHub release is live, verify install/upgrade from the tap when practical:

```sh
brew update
brew reinstall ossianhempel/tap/things3-cli
things --version
```

Do not claim Homebrew is ready until the formula references the released tag and tarball checksums.

## Publishing

Publishing is automated. Do **not** tag by hand and do **not** run the local
publish path for normal releases. Instead:

1. Land the changelog + formula bump on `main` via a release PR.
2. On merge, `.github/workflows/release.yml` detects the new `## [X.Y.Z]`
   CHANGELOG section, tags `vX.Y.Z` at the merge commit, runs tests, builds
   artifacts, generates notes, and publishes `things3-cli vX.Y.Z`. The workflow
   is idempotent: if the tag already exists it skips.
3. Manual override only if needed: Actions → Release → Run workflow, enter the
   version (`workflow_dispatch`).

The legacy `git tag … && git push origin vX.Y.Z` and `./scripts/release.sh`
paths still work but are not the standard flow.

## Homebrew Tap (agent-driven, after the release is live)

CI does not touch the tap. Once the GitHub release exists (so the released
tarballs are downloadable and their checksums are final), the agent updates and
pushes the tap:

```sh
./scripts/update-brew-formula.sh --version vX.Y.Z --tap-dir ~/Developer/homebrew-tap
git -C ~/Developer/homebrew-tap add Formula/things3-cli.rb
git -C ~/Developer/homebrew-tap commit -m "things3-cli vX.Y.Z"
git -C ~/Developer/homebrew-tap push
```

The maintainer never does this manually — the agent driving the release does.

## Post-Release Verification

Verify:

- GitHub release title is `things3-cli vX.Y.Z`.
- Release body contains only that version's changelog bullets.
- Assets include both darwin tarballs and `checksums.txt`.
- Formula version, URLs, and checksums point to the released version.
- Installed `things --version` reports `vX.Y.Z`.

## Agent Skill Sync

If command behavior changed, update:

- `skills/things/SKILL.md`
- `../agent-scripts/archived-skills/things/SKILL.md` when that mirror exists

If release mechanics changed, update:

- `.codex/skills/release-flow/SKILL.md`
- `.codex/skills/release-flow/references/release-flow.md`
- `docs/RELEASING.md`
