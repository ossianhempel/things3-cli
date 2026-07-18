#!/bin/sh
set -eu

script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd -P)
repo_root=$(CDPATH= cd -- "$script_dir/.." && pwd -P)
canonical="$repo_root/skills/things/SKILL.md"
mirror_root=${THINGS_SKILL_MIRROR_ROOT:-"$repo_root/../agent-scripts"}
expected_remote='ossianhempel/agent-scripts'
destination_rel='skills/things/SKILL.md'
destination="$mirror_root/$destination_rel"

if [ ! -f "$canonical" ] || [ -L "$canonical" ]; then
	printf 'refusing sync: canonical skill is missing, not regular, or symlinked: %s\n' "$canonical" >&2
	exit 1
fi
if [ -L "$mirror_root" ] || [ ! -d "$mirror_root" ] || { [ ! -d "$mirror_root/.git" ] && [ ! -f "$mirror_root/.git" ]; }; then
	printf 'refusing sync: agent-scripts repository is missing: %s\n' "$mirror_root" >&2
	exit 1
fi

remote=$(git -C "$mirror_root" remote get-url origin 2>/dev/null || true)
case "$remote" in
	"https://github.com/$expected_remote"|\
	"https://github.com/$expected_remote.git"|\
	"ssh://git@github.com/$expected_remote"|\
	"ssh://git@github.com/$expected_remote.git"|\
	"git@github.com:$expected_remote"|\
	"git@github.com:$expected_remote.git") ;;
	*)
		printf 'refusing sync: %s is not the expected %s repository (origin: %s)\n' "$mirror_root" "$expected_remote" "${remote:-missing}" >&2
		exit 1
		;;
esac

# Refuse symlinks anywhere in the fixed destination chain. This check occurs
# before creating directories or temporary files, so failures cannot mutate it.
path="$mirror_root"
for component in skills things SKILL.md; do
	path="$path/$component"
	if [ -L "$path" ]; then
		printf 'refusing sync: destination chain contains symlink: %s\n' "$path" >&2
		exit 1
	fi
done
if [ ! -d "$mirror_root/skills/things" ]; then
	printf 'refusing sync: exact active skill directory is missing: %s\n' "$mirror_root/skills/things" >&2
	exit 1
fi
if [ -e "$destination" ] && [ ! -f "$destination" ]; then
	printf 'refusing sync: destination is not a regular file: %s\n' "$destination" >&2
	exit 1
fi

tmp=$(mktemp "$mirror_root/skills/things/.SKILL.md.sync.XXXXXX")
trap 'rm -f -- "$tmp"' EXIT HUP INT TERM
cp "$canonical" "$tmp"
chmod --reference="$canonical" "$tmp" 2>/dev/null || chmod 644 "$tmp"
mv -f "$tmp" "$destination"
trap - EXIT HUP INT TERM
printf 'synchronized canonical Things skill to %s\n' "$destination"
