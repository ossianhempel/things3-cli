#!/bin/sh
set -eu

script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd -P)
repo_root=$(CDPATH= cd -- "$script_dir/.." && pwd -P)
canonical="$repo_root/skills/things/SKILL.md"
mirror_root=${THINGS_SKILL_MIRROR_ROOT:-"$repo_root/../agent-scripts"}
mirror="$mirror_root/skills/things/SKILL.md"
require_mirror=false

if [ "${1:-}" = "--require-mirror" ]; then
	require_mirror=true
elif [ "$#" -ne 0 ]; then
	printf 'usage: %s [--require-mirror]\n' "$0" >&2
	exit 2
fi

if [ ! -f "$canonical" ] || [ -L "$canonical" ]; then
	printf 'canonical Things skill is missing, not regular, or symlinked: %s\n' "$canonical" >&2
	exit 1
fi

for required in \
	'## Safe operating workflow' \
	'**Identify**' \
	'**Preview**' \
	'**Verify**' \
	'Repeating projects are unsupported' \
	'partial success' \
	'Full Disk Access'
do
	if ! grep -Fq -- "$required" "$canonical"; then
		printf 'canonical Things skill is missing required guidance: %s\n' "$required" >&2
		exit 1
	fi
done

if [ ! -d "$mirror_root/.git" ] && [ ! -f "$mirror_root/.git" ]; then
	if [ "$require_mirror" = true ]; then
		printf 'agent-scripts checkout is required but missing: %s\n' "$mirror_root" >&2
		exit 1
	fi
	printf 'canonical Things skill validated; sibling agent-scripts checkout not present, parity skipped\n'
	exit 0
fi

if [ ! -f "$mirror" ] || [ -L "$mirror" ]; then
	printf 'active Things skill mirror is missing, not regular, or symlinked: %s\n' "$mirror" >&2
	exit 1
fi

if ! cmp -s "$canonical" "$mirror"; then
	printf 'Things skill drift detected:\n  canonical: %s\n  mirror:    %s\n' "$canonical" "$mirror" >&2
	diff -u "$mirror" "$canonical" >&2 || true
	exit 1
fi

printf 'Things skill mirror matches canonical: %s\n' "$mirror"
