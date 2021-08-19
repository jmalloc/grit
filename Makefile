-include .makefiles/Makefile
-include .makefiles/pkg/go/v1/Makefile

run: $(GO_DEBUG_DIR)/grit2
	$< $(RUN_ARGS)

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"

