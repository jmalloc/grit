SHELL := /bin/bash
-include artifacts/make/go.mk

artifacts/make/%.mk:
	bash <(curl -s https://rinq.github.io/make/install) $*
