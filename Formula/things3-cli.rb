class Things3Cli < Formula
  desc "CLI for Things 3"
  homepage "https://github.com/ossianhempel/things3-cli"
  version "0.4.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/ossianhempel/things3-cli/releases/download/v0.4.0/things-0.4.0-darwin-arm64.tar.gz"
      sha256 "0c48382c37ee5f9f783428be5ec1def9e1340dea0669b818b038e81e257a76b8"
    else
      url "https://github.com/ossianhempel/things3-cli/releases/download/v0.4.0/things-0.4.0-darwin-amd64.tar.gz"
      sha256 "a9bae79196f894c29b790108eb1066b3997df09f246cd1892761ed64c3ce1856"
    end
  end

  def install
    bin.install "things"
  end

  test do
    system "#{bin}/things", "--version"
  end
end
