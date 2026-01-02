class Things3Cli < Formula
  desc "CLI for Things 3"
  homepage "https://github.com/ossianhempel/things3-cli"
  url "https://github.com/ossianhempel/things3-cli/archive/959aeea72c81cf93a91c9811dd7c127e9b21b7d8.tar.gz"
  sha256 "a994a73e258b97a4e941c636613c295a0a36061f7e7e205cf7930777e58aadad"
  version "20260102103314"

  depends_on "go" => :build

  def install
    ld_version = "959aeea"
    ldflags = "-s -w -X github.com/ossianhempel/things3-cli/internal/cli.Version=#{ld_version}"
    system "go", "build", "-trimpath", "-ldflags", ldflags, "-o", bin/"things", "./cmd/things"
  end

  test do
    system "#{bin}/things", "--version"
  end
end
