#!/usr/bin/env bash
set -e
set -o pipefail

TAG="$1"
HASH_AMD64="$(shasum -a 256 artifacts/archives/grit-${TAG}-darwin-amd64.zip | awk '{ print $1 }')"
HASH_ARM64="$(shasum -a 256 artifacts/archives/grit-${TAG}-darwin-arm64.zip | awk '{ print $1 }')"

dir=$(mktemp -d)
git clone "https://${GITHUB_TOKEN}@github.com/jmalloc/homebrew-grit" "$dir"
cd "$dir"

tee grit.rb <<EOF
class Grit < Formula
  desc "Keep track of your local Git clones."
  homepage "https://github.com/jmalloc/grit"

  version "${TAG}"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/jmalloc/grit/releases/download/${TAG}/grit-${TAG}-darwin-amd64.zip"
      sha256 "${HASH_AMD64}"
    end

    if Hardware::CPU.arm?
      url "https://github.com/jmalloc/grit/releases/download/${TAG}/grit-${TAG}-darwin-arm64.zip"
      sha256 "${HASH_ARM64}"
    end
  end


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
