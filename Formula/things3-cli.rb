class Things3Cli < Formula
  desc "CLI for Things 3"
  homepage "https://github.com/ossianhempel/things3-cli"
  version "0.1.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/ossianhempel/things3-cli/releases/download/v0.1.0/things-0.1.0-darwin-arm64.tar.gz"
      sha256 "db92fc7178a32ff338222631c7056f37690942a462f6cc876c7c66efd443f600"
    else
      url "https://github.com/ossianhempel/things3-cli/releases/download/v0.1.0/things-0.1.0-darwin-amd64.tar.gz"
      sha256 "f40774fb313f7d840aa69f0519f69f3963736eb7b83455692b71bedf815d5747"
    end
  end

  def install
    bin.install "things"
  end

  test do
    system "#{bin}/things", "--version"
  end
end
