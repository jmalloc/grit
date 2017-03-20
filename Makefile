SHELL := /bin/bash
-include artifacts/make/go.mk

run: artifacts/build/debug/$(GOOS)/$(GOARCH)/grit
	GRIT_CONFIG=etc/testing.toml "$<" $(RUN_ARGS)

artifacts/make/%.mk:
	bash <(curl -s https://rinq.github.io/make/install) $*
