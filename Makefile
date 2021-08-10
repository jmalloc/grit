GO_MATRIX_OS := darwin linux windows

-include .makefiles/Makefile
-include .makefiles/pkg/go/v1/Makefile
-include .makefiles/ext/na4ma4/lib/goreleaser/v1/Makefile

run: artifacts/build/debug/$(GOHOSTOS)/$(GOHOSTARCH)/grit
	GRIT_CONFIG=etc/testing.toml "$<" $(RUN_ARGS)

# homebrew: artifacts/archives/grit-$(GIT_HEAD_TAG)-darwin-amd64.zip
# 	bin/homebrew.sh "$(GIT_HEAD_TAG)" "$<"

.makefiles/ext/na4ma4/%: .makefiles/Makefile
	@curl -sfL https://raw.githubusercontent.com/na4ma4/makefiles-ext/main/v1/install | bash /dev/stdin "$@"

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"
