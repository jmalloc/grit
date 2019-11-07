MATRIX_OS := darwin linux windows

-include .makefiles/Makefile
-include .makefiles/pkg/go/v1/Makefile

# run: artifacts/build/debug/$(GOOS)/$(GOARCH)/grit
# 	GRIT_CONFIG=etc/testing.toml "$<" $(RUN_ARGS)

# homebrew: artifacts/archives/grit-darwin-amd64.tar.gz
# 	bin/homebrew.sh "$<"

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"
