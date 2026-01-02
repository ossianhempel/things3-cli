---
summary: 'things3-cli release checklist'
read_when:
  - preparing a release
  - writing release notes
---

# Releasing things3-cli

## Guardrails

- Title every GitHub release as `things3-cli <version>`.
- Release body = the CHANGELOG bullets for that version only (no extra prose).

## Checklist

- [ ] Update `CHANGELOG.md`: move items from **Unreleased** into a new version section with today’s date.
- [ ] Run tests: `make test`.
- [ ] Build release artifacts: `./scripts/build-release.sh vX.Y.Z`.
- [ ] Update Homebrew formula:
  `./scripts/update-brew-formula.sh --version vX.Y.Z --tap-dir ~/Developer/homebrew-tap`
  (commits to `Formula/` and copies into the tap repo).
- [ ] Generate release notes: `./scripts/release-notes.sh vX.Y.Z` (should output only the changelog bullets).
- [ ] Tag the release: `git tag vX.Y.Z` then `git push --tags`.
- [ ] Create the GitHub release:
  - Option A (local): `./scripts/release.sh vX.Y.Z`
  - Option B (CI): push the tag and let `.github/workflows/release.yml` publish the release
- [ ] Verify the GitHub release title, notes, and assets match the guardrails.
