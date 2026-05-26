---
name: release-flow
description: Maintainer release workflow for things3-cli. Use when preparing, validating, tagging, publishing, or verifying a public things3-cli release or Homebrew formula update.
---

# release-flow

Use this skill when shipping `things3-cli` to external users and agents.

Start here:
- Read `docs/RELEASING.md`.
- Read `references/release-flow.md` for the validation sequence and command details.

Core rule:
- Do not publish a release from assumptions. Verify tests, artifacts, changelog notes, formula checksums, GitHub release assets, and the installed binary path before calling the release done.

Tagging and the GitHub release are automated: merging a PR that adds a new
`## [X.Y.Z]` section to `CHANGELOG.md` makes `.github/workflows/release.yml` tag
and publish `vX.Y.Z`. Never tag by hand. The Homebrew tap is NOT in CI — the
agent updates and pushes it after the release is live.

Quick path:
1. Confirm the release version (SemVer scope of the merged changes).
2. Add a `## [X.Y.Z] - YYYY-MM-DD` section to `CHANGELOG.md` with the bullets.
3. Build artifacts (`./scripts/build-release.sh vX.Y.Z`) and bump the in-repo
   formula (`./scripts/update-brew-formula.sh --version vX.Y.Z`) so it ships in the PR.
4. Run `make test`; smoke-test a built binary's `--version`.
5. Sanity-check release notes with `./scripts/release-notes.sh vX.Y.Z`.
6. Open a release PR and merge it — CI tags and publishes the release.
7. After CI publishes, verify the GitHub release, then update and push the tap:
   `./scripts/update-brew-formula.sh --version vX.Y.Z --tap-dir ~/Developer/homebrew-tap`
   and commit + push `~/Developer/homebrew-tap`.
8. Verify the Homebrew install path.

When CLI behavior changes, also update `skills/things/SKILL.md` and the mirrored `../agent-scripts/archived-skills/things/SKILL.md` if that mirror exists.
