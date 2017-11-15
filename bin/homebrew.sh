set -e
set -o pipefail

dir=$(mktemp -d)
git clone "https://${GITHUB_TOKEN}@github.com/jmalloc/homebrew-grit" "$dir"
cd "$dir"

tee grit.rb <<EOF
class Grit < Formula
  desc "Keep track of your local Git clones."
  homepage "https://github.com/jmalloc/grit"

  version "${TRAVIS_TAG}"
  url "https://github.com/jmalloc/grit/releases/download/${TRAVIS_TAG}/grit-darwin-amd64.tar.gz"
  sha256 "$(shasum -a 256 "$1" | awk '{ print $1 }')"

  def install
      bin.install "grit"
  end

  test do
    system "grit"
  end
end
EOF

git commit -a -m 'Update to v${TRAVIS_TAG}'
