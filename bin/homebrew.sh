#!/usr/bin/env bash
set -e
set -o pipefail

TAG="$1"
HASH="$(shasum -a 256 $2 | awk '{ print $1 }')"

dir=$(mktemp -d)
git clone "https://${GITHUB_TOKEN}@github.com/jmalloc/homebrew-grit" "$dir"
cd "$dir"

tee grit.rb <<EOF
class Grit < Formula
  desc "Keep track of your local Git clones."
  homepage "https://github.com/jmalloc/grit"

  version "${1}"
  url "https://github.com/jmalloc/grit/releases/download/${TAG}/grit-${TAG}-darwin-amd64.zip"
  sha256 "${HASH}"

  def install
      bin.install "grit"
  end

  test do
    system "grit"
  end
end
EOF

git commit -a -m "Update to v${TAG}"
git push
