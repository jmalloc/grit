GO_MATRIX += darwin/amd64 darwin/arm64
GO_MATRIX += linux/amd64
GO_MATRIX += windows/amd64

-include .makefiles/Makefile
-include .makefiles/pkg/go/v1/Makefile

run: $(GO_DEBUG_DIR)/grit
	GRIT_CONFIG=etc/testing.toml "$<" $(RUN_ARGS)

homebrew: artifacts/archives/grit-$(GIT_HEAD_TAG)-darwin-amd64.zip artifacts/archives/grit-$(GIT_HEAD_TAG)-darwin-arm64.zip
	bin/homebrew.sh "$(GIT_HEAD_TAG)"

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"
