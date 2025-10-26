class TrelloCli < Formula
  desc "A comprehensive Trello CLI tool optimized for LLM integration"
  homepage "https://github.com/danbruder/trello-cli"
  url "https://github.com/danbruder/trello-cli/archive/v1.0.0.tar.gz"
  sha256 ""
  license "MIT"
  head "https://github.com/danbruder/trello-cli.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", "-ldflags", "-s -w", "-o", "trello-cli", "."
    bin.install "trello-cli"
  end

  test do
    # Test that the binary works
    system "#{bin}/trello-cli", "--help"
  end
end
