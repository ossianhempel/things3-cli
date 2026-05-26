---
summary: 'things3-cli release checklist'
read_when:
  - preparing a release
  - writing release notes
---

# Releasing things3-cli

## How releases ship

Tagging and the GitHub release are **automated**. `.github/workflows/release.yml`
watches `main`: when a merged commit adds a new `## [X.Y.Z]` section to
`CHANGELOG.md`, CI tags `vX.Y.Z` (at the merge commit), runs tests, builds the
darwin tarballs, generates the release notes from that changelog section, and
publishes `things3-cli vX.Y.Z`. The workflow is idempotent — a later CHANGELOG
edit that only touches **Unreleased** finds the tag already present and skips.
A `workflow_dispatch` (Actions → Release → Run workflow, enter the version)
is the manual override.

So nobody tags by hand. The maintainer's only required action is to **merge the
release PR**. Updating the Homebrew tap is **not** in CI — the agent driving the
release does it (see below).

## Guardrails

- Title every GitHub release as `things3-cli <version>` (CI does this).
- Release body = the CHANGELOG bullets for that version only (no extra prose).
- Do not call a release complete until tests, artifacts, release notes,
  Homebrew formula checksums, GitHub release assets, and the installed binary
  version have been verified.

## Checklist (the release PR)

- [ ] Confirm the next version matches SemVer scope for the merged changes.
- [ ] Update `CHANGELOG.md`: add a new `## [X.Y.Z] - YYYY-MM-DD` section with
      the bullets for this release.
- [ ] Bump the in-repo Homebrew formula so it ships in the same PR:
  `./scripts/build-release.sh vX.Y.Z` then
  `./scripts/update-brew-formula.sh --version vX.Y.Z` (writes `Formula/things3-cli.rb`).
- [ ] Run tests: `make test`.
- [ ] Smoke-test an artifact binary and confirm `things --version` reports `vX.Y.Z`.
- [ ] Sanity-check release notes: `./scripts/release-notes.sh vX.Y.Z`
      (should output only the changelog bullets).
- [ ] Open the PR and merge it. CI tags and publishes the release automatically.

## Checklist (after CI publishes — the agent does this, not the maintainer)

- [ ] Confirm the GitHub release `things3-cli vX.Y.Z` is live with both darwin
      tarballs and `checksums.txt` attached, and the body is the changelog bullets.
- [ ] Update and push the Homebrew tap (the released tarballs must already exist
      so the checksums resolve):
  `./scripts/update-brew-formula.sh --version vX.Y.Z --tap-dir ~/Developer/homebrew-tap`
  then commit and push the tap repo.
- [ ] Verify the tap formula points at the new release URLs and checksums.
- [ ] Verify Homebrew install or reinstall from `ossianhempel/tap/things3-cli`
      and confirm `things --version`.
