class Things3Cli < Formula
  desc "CLI for Things 3"
  homepage "https://github.com/ossianhempel/things3-cli"
  version "0.2.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/ossianhempel/things3-cli/releases/download/v0.2.0/things-0.2.0-darwin-arm64.tar.gz"
      sha256 "e6b84fee29e83cb39dc4b6a72480a6de05fc5876723dcc1f0a8e4e7c5f7f783f"
    else
      url "https://github.com/ossianhempel/things3-cli/releases/download/v0.2.0/things-0.2.0-darwin-amd64.tar.gz"
      sha256 "ce9e1dadf7198c5c2f51d5eaf94f925110cf0170565e4066624802c2b104ca72"
    end
  end

  def install
    bin.install "things"
  end

  test do
    system "#{bin}/things", "--version"
  end
end
