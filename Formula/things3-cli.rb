class Things3Cli < Formula
  desc "CLI for Things 3"
  homepage "https://github.com/ossianhempel/things3-cli"
  version "0.5.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/ossianhempel/things3-cli/releases/download/v0.5.0/things-0.5.0-darwin-arm64.tar.gz"
      sha256 "253f31237748fd8282e825567829a1460e90d94e4209554c92ae689781409bc2"
    else
      url "https://github.com/ossianhempel/things3-cli/releases/download/v0.5.0/things-0.5.0-darwin-amd64.tar.gz"
      sha256 "bc3c245b28697dbee22e77ebcfc2c8e483bc14e5bd481ed6f5b1bac558c2e39e"
    end
  end

  def install
    bin.install "things"
  end

  test do
    system "#{bin}/things", "--version"
  end
end
