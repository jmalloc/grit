GO_EMBEDDED_FILES += $(shell find cmd -iname '*.txt')

-include .makefiles/Makefile
-include .makefiles/pkg/go/v1/Makefile

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"

run: $(GO_DEBUG_DIR)/grit2
	$< $(RUN_ARGS)
