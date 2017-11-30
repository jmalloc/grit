MATRIX_OS := darwin linux windows

-include artifacts/make/go/Makefile

run: artifacts/build/debug/$(GOOS)/$(GOARCH)/grit
	GRIT_CONFIG=etc/testing.toml "$<" $(RUN_ARGS)

homebrew: artifacts/archives/grit-darwin-amd64.tar.gz
	bin/homebrew.sh "$<"

artifacts/make/%/Makefile:
	curl -sf https://jmalloc.github.io/makefiles/fetch | bash /dev/stdin $*
