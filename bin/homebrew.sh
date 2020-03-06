set -e
set -o pipefail

hash="$(shasum -a 256 $1 | awk '{ print $1 }')"
tag="${GITHUB_REF#refs/tags/}"
dir=$(mktemp -d)
git clone "https://${GITHUB_TOKEN}@github.com/jmalloc/homebrew-grit" "$dir"
cd "$dir"

tee grit.rb <<EOF
class Grit < Formula
  desc "Keep track of your local Git clones."
  homepage "https://github.com/jmalloc/grit"

  version "${tag}"
  url "https://github.com/jmalloc/grit/releases/download/${tag}/grit-${tag}-darwin-amd64.zip"
  sha256 "${hash}"

  def install
      bin.install "grit"
  end

  test do
    system "grit"
  end
end
EOF

git commit -a -m "Update to v${tag}"
git push
