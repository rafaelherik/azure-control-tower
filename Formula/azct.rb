class Azct < Formula
  desc "Terminal-based UI for exploring and managing Azure resources"
  homepage "https://github.com/rafaelherik/azure-control-tower"
  url "https://github.com/rafaelherik/azure-control-tower/archive/v0.0.1.tar.gz"
  sha256 ""
  license "MIT"
  head "https://github.com/rafaelherik/azure-control-tower.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", "-ldflags", "-s -w -X main.version=#{version}", "-o", bin/"azct", "./cmd/azct"
  end

  test do
    system "#{bin}/azct", "--version"
  end
end

