#!/bin/sh
set -eu

script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd -P)
repo_root=$(CDPATH= cd -- "$script_dir/.." && pwd -P)
test_root=$(mktemp -d "${TMPDIR:-/tmp}/things-skill-sync-test.XXXXXX")
trap 'rm -rf -- "$test_root"' EXIT HUP INT TERM

fixture_root="$test_root/things3-cli"
mkdir -p "$fixture_root/scripts" "$fixture_root/skills/things"
cp "$repo_root/scripts/sync-things-skill.sh" "$fixture_root/scripts/"
printf 'canonical test skill\n' >"$fixture_root/skills/things/SKILL.md"
chmod +x "$fixture_root/scripts/sync-things-skill.sh"

new_mirror() {
	name=$1
	remote=$2
	mirror="$test_root/$name"
	git init -q "$mirror"
	git -C "$mirror" remote add origin "$remote"
	mkdir -p "$mirror/skills/things"
	printf 'old mirror\n' >"$mirror/skills/things/SKILL.md"
	printf '%s\n' "$mirror"
}

expect_failure() {
	label=$1
	shift
	if "$@" >"$test_root/$label.out" 2>"$test_root/$label.err"; then
		printf 'expected failure: %s\n' "$label" >&2
		exit 1
	fi
}

# Repository identity must match the normalized GitHub origin exactly. A host
# whose name merely ends in github.com must never be accepted.
mirror=$(new_mirror identity 'https://evilgithub.com/ossianhempel/agent-scripts.git')
expect_failure identity env THINGS_SKILL_MIRROR_ROOT="$mirror" "$fixture_root/scripts/sync-things-skill.sh"
grep -Fq 'not the expected ossianhempel/agent-scripts repository' "$test_root/identity.err"
grep -Fq 'old mirror' "$mirror/skills/things/SKILL.md"

# Refuse a symlink at any component of the fixed destination chain.
mirror=$(new_mirror symlink 'git@github.com:ossianhempel/agent-scripts.git')
rm -rf "$mirror/skills"
mkdir "$test_root/symlink-target"
ln -s "$test_root/symlink-target" "$mirror/skills"
expect_failure symlink env THINGS_SKILL_MIRROR_ROOT="$mirror" "$fixture_root/scripts/sync-things-skill.sh"
grep -Fq 'destination chain contains symlink' "$test_root/symlink.err"

# Refuse an existing destination that is not a regular file.
mirror=$(new_mirror nonregular 'ssh://git@github.com/ossianhempel/agent-scripts')
rm "$mirror/skills/things/SKILL.md"
mkdir "$mirror/skills/things/SKILL.md"
expect_failure nonregular env THINGS_SKILL_MIRROR_ROOT="$mirror" "$fixture_root/scripts/sync-things-skill.sh"
grep -Fq 'destination is not a regular file' "$test_root/nonregular.err"

# A valid checkout is updated byte-for-byte without leaving a temporary file.
mirror=$(new_mirror success 'https://github.com/ossianhempel/agent-scripts.git')
env THINGS_SKILL_MIRROR_ROOT="$mirror" "$fixture_root/scripts/sync-things-skill.sh" >"$test_root/success.out"
cmp -s "$fixture_root/skills/things/SKILL.md" "$mirror/skills/things/SKILL.md"
if find "$mirror/skills/things" -name '.SKILL.md.sync.*' -print -quit | grep -q .; then
	printf 'successful sync left a temporary file\n' >&2
	exit 1
fi

# If the atomic rename fails, preserve the old destination and clean the
# temporary file through the EXIT trap.
mirror=$(new_mirror atomic 'git@github.com:ossianhempel/agent-scripts')
fake_bin="$test_root/fake-bin"
mkdir "$fake_bin"
cat >"$fake_bin/mv" <<'EOF'
#!/bin/sh
exit 73
EOF
chmod +x "$fake_bin/mv"
expect_failure atomic env PATH="$fake_bin:$PATH" THINGS_SKILL_MIRROR_ROOT="$mirror" "$fixture_root/scripts/sync-things-skill.sh"
grep -Fq 'old mirror' "$mirror/skills/things/SKILL.md"
if find "$mirror/skills/things" -name '.SKILL.md.sync.*' -print -quit | grep -q .; then
	printf 'failed sync left a temporary file\n' >&2
	exit 1
fi

printf 'Things skill sync safety tests passed\n'
