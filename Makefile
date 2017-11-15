MATRIX_OS := darwin linux windows
SHELL := /bin/bash
-include artifacts/make/go.mk

run: artifacts/build/debug/$(GOOS)/$(GOARCH)/grit
	GRIT_CONFIG=etc/testing.toml "$<" $(RUN_ARGS)

artifacts/make/%.mk:
	bash <(curl -s https://rinq.github.io/make/install) $*

homebrew: artifacts/archives/grit-darwin-amd64.tar.gz
	bin/homebrew.sh "$<"
