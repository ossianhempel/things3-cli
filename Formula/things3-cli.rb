class Things3Cli < Formula
  desc "CLI for Things 3"
  homepage "https://github.com/ossianhempel/things3-cli"
  version "0.1.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/ossianhempel/things3-cli/releases/download/v0.1.0/things-0.1.0-darwin-arm64.tar.gz"
      sha256 "ee98aa7e5656475bbe7598251576529db2f45a15e292d20b3fefa12e867889d4"
    else
      url "https://github.com/ossianhempel/things3-cli/releases/download/v0.1.0/things-0.1.0-darwin-amd64.tar.gz"
      sha256 "69cc7310b88159ed480967019a9ba1d0341ab54f07692f4c36a29903a84ad129"
    end
  end

  def install
    bin.install "things"
  end

  test do
    system "#{bin}/things", "--version"
  end
end
