class Things3Cli < Formula
  desc "CLI for Things 3"
  homepage "https://github.com/ossianhempel/things3-cli"
  url "https://github.com/ossianhempel/things3-cli/archive/d2c4b46a400656ff35ff249aa58ae0ecd5e168e2.tar.gz"
  sha256 "238bc0e6e278b49010f57b66d2329d00609c177d094d278c251b856da9f7d8c7"
  version "d2c4b46"

  depends_on "go" => :build

  def install
    ldflags = "-s -w -X github.com/ossianhempel/things3-cli/internal/cli.Version=#{version}"
    system "go", "build", "-trimpath", "-ldflags", ldflags, "-o", bin/"things", "./cmd/things"
  end

  test do
    system "#{bin}/things", "--version"
  end
end
