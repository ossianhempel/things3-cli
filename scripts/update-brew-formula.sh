#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$ROOT_DIR"

VERSION=""
TAP_DIR=""
while [[ $# -gt 0 ]]; do
  case "$1" in
    --version)
      VERSION=${2:-}
      if [[ -z "$VERSION" ]]; then
        echo "ERROR: --version requires a value (e.g. v0.1.0)." >&2
        exit 1
      fi
      shift 2
      ;;
    --tap-dir)
      TAP_DIR=${2:-}
      if [[ -z "$TAP_DIR" ]]; then
        echo "ERROR: --tap-dir requires a path." >&2
        exit 1
      fi
      shift 2
      ;;
    *)
      echo "usage: scripts/update-brew-formula.sh [--version vX.Y.Z] [--tap-dir PATH]" >&2
      exit 1
      ;;
  esac
done

mkdir -p "$ROOT_DIR/Formula"

if [[ -n "$VERSION" ]]; then
  VERSION=${VERSION#v}
  TAG="v${VERSION}"
  CHECKSUMS="$ROOT_DIR/dist/checksums.txt"

  if [[ ! -f "$CHECKSUMS" ]]; then
    echo "ERROR: $CHECKSUMS not found. Run scripts/build-release.sh $TAG first." >&2
    exit 1
  fi

  ARM_TARBALL="things-${VERSION}-darwin-arm64.tar.gz"
  AMD_TARBALL="things-${VERSION}-darwin-amd64.tar.gz"
  ARM_SHA=$(awk -v f="$ARM_TARBALL" '$2==f {print $1}' "$CHECKSUMS")
  AMD_SHA=$(awk -v f="$AMD_TARBALL" '$2==f {print $1}' "$CHECKSUMS")

  if [[ -z "$ARM_SHA" || -z "$AMD_SHA" ]]; then
    echo "ERROR: Missing checksums for release tarballs in $CHECKSUMS." >&2
    exit 1
  fi

  cat > "$ROOT_DIR/Formula/things3-cli.rb" <<FORMULA
class Things3Cli < Formula
  desc "CLI for Things 3"
  homepage "https://github.com/ossianhempel/things3-cli"
  version "${VERSION}"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/ossianhempel/things3-cli/releases/download/${TAG}/${ARM_TARBALL}"
      sha256 "${ARM_SHA}"
    else
      url "https://github.com/ossianhempel/things3-cli/releases/download/${TAG}/${AMD_TARBALL}"
      sha256 "${AMD_SHA}"
    end
  end

  def install
    bin.install "things"
  end

  test do
    system "#{bin}/things", "--version"
  end
end
FORMULA

  if [[ -n "$TAP_DIR" ]]; then
    mkdir -p "$TAP_DIR/Formula"
    cp "$ROOT_DIR/Formula/things3-cli.rb" "$TAP_DIR/Formula/things3-cli.rb"
    echo "Copied formula to $TAP_DIR/Formula/things3-cli.rb"
  fi

  echo "Wrote Formula/things3-cli.rb for release ${TAG}"
  exit 0
fi

COMMIT=$(git rev-parse HEAD)
SHORT_COMMIT=${COMMIT:0:7}
URL="https://github.com/ossianhempel/things3-cli/archive/${COMMIT}.tar.gz"

if ! command -v curl >/dev/null 2>&1; then
  echo "ERROR: curl is required." >&2
  exit 1
fi

SHA=$(curl -fsSL "$URL" | shasum -a 256 | awk '{print $1}')

  cat > "$ROOT_DIR/Formula/things3-cli.rb" <<FORMULA
class Things3Cli < Formula
  desc "CLI for Things 3"
  homepage "https://github.com/ossianhempel/things3-cli"
  url "${URL}"
  sha256 "${SHA}"
  version "${SHORT_COMMIT}"

  depends_on "go" => :build

  def install
    ldflags = "-s -w -X github.com/ossianhempel/things3-cli/internal/cli.Version=#{version}"
    system "go", "build", "-trimpath", "-ldflags", ldflags, "-o", bin/"things", "./cmd/things"
  end

  test do
    system "#{bin}/things", "--version"
  end
end
FORMULA

if [[ -n "$TAP_DIR" ]]; then
  mkdir -p "$TAP_DIR/Formula"
  cp "$ROOT_DIR/Formula/things3-cli.rb" "$TAP_DIR/Formula/things3-cli.rb"
  echo "Copied formula to $TAP_DIR/Formula/things3-cli.rb"
fi

echo "Wrote Formula/things3-cli.rb for commit ${COMMIT}"
