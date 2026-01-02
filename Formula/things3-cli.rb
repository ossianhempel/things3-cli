class Things3Cli < Formula
  desc "CLI for Things 3"
  homepage "https://github.com/ossianhempel/things3-cli"
  url "https://github.com/ossianhempel/things3-cli/archive/db9b7deb37876a0b06d7aefb5d28292df7183a4c.tar.gz"
  sha256 "d3236a2886893417ee2d9421248b50daf8d728deef4e2e2e8cb3b7a9d50c409c"
  version "db9b7de"

  depends_on "go" => :build

  def install
    ldflags = "-s -w -X github.com/ossianhempel/things3-cli/internal/cli.Version=#{version}"
    system "go", "build", "-trimpath", "-ldflags", ldflags, "-o", bin/"things", "./cmd/things"
  end

  test do
    system "#{bin}/things", "--version"
  end
end
