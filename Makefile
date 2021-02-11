GO_MATRIX_OS := darwin linux windows

-include .makefiles/Makefile
-include .makefiles/pkg/go/v1/Makefile

run: artifacts/build/debug/$(GOHOSTOS)/$(GOHOSTARCH)/grit
	GRIT_CONFIG=etc/testing.toml "$<" $(RUN_ARGS)

homebrew: artifacts/archives/grit-$(GIT_HEAD_TAG)-darwin-amd64.zip
	bin/homebrew.sh "$(GIT_HEAD_TAG)" "$<"

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"

SHELL := /bin/bash -o pipefail

clean:: clean-test

.PHONY: test-gritrepos
test-gritrepos:
	-$(RM) -r /tmp/grit-repos/
	mkdir -p /tmp/grit-repos/
	mkdir -p /tmp/grit-repos/testorg/testrepo.git
	cd /tmp/grit-repos/testorg/testrepo.git && git init && touch file.dat && git add -A . && git commit -am "initial commit"

## Testing empty repo clone
TEST_TARGETS += /tmp/grit-test/clone/github.com/git-fixtures/empty/.git
/tmp/grit-test/clone/github.com/git-fixtures/empty/.git: artifacts/build/debug/$(GOHOSTOS)/$(GOHOSTARCH)/grit
	-$(RM) -r "$(@D)"
	GRIT_CONFIG=etc/testing.toml "$(<)" clone git-fixtures/empty
	test -d "$(@)"
	cd "$(@D)" && git status

## Testing malformed slug
TEST_TARGETS += artifacts/test/malformed-slug.txt
artifacts/test/malformed-slug.txt: artifacts/build/debug/$(GOHOSTOS)/$(GOHOSTARCH)/grit
	-@mkdir -p "$(@D)"
	-GRIT_CONFIG=etc/testing-file.toml "$(<)" clone sluggyslug
	test ! -d "/tmp/grit-repos/sluggyslug"

## Testing normal clone
TEST_TARGETS += artifacts/test/normal-clone.txt
artifacts/test/normal-clone.txt: artifacts/build/debug/$(GOHOSTOS)/$(GOHOSTARCH)/grit test-gritrepos
	-@mkdir -p "$(@D)"
	GRIT_CONFIG=etc/testing-file.toml "$(<)" clone testorg/testrepo
	test -d "/tmp/grit-test/clone/tmp/grit-repos/testorg/testrepo/.git"
	ls -al "/tmp/grit-test/clone/tmp/grit-repos/testorg/testrepo"

## Testing empty clone index
TEST_TARGETS += artifacts/test/empty-clone-index.txt
artifacts/test/empty-clone-index.txt: artifacts/build/debug/$(GOHOSTOS)/$(GOHOSTARCH)/grit artifacts/test/normal-clone.txt
	-@mkdir -p "$(@D)"
	GRIT_CONFIG=etc/testing-file.toml "$(<)" index ls | grep empty | tee "$(@)"
	GRIT_CONFIG=etc/testing-file.toml "$(<)" index ls | grep git-fixtures/empty | tee -a "$(@)"

## Testing normal clone index
TEST_TARGETS += artifacts/test/normal-clone-index.txt
artifacts/test/normal-clone-index.txt: artifacts/build/debug/$(GOHOSTOS)/$(GOHOSTARCH)/grit artifacts/test/normal-clone.txt
	-@mkdir -p "$(@D)"
	GRIT_CONFIG=etc/testing-file.toml "$(<)" index ls | grep testrepo | tee "$(@)"
	GRIT_CONFIG=etc/testing-file.toml "$(<)" index ls | grep tmp/grit-repos/testorg/testrepo | tee -a "$(@)"

## Testing empty clone change directory
TEST_TARGETS += artifacts/test/empty-clone-cd.txt
artifacts/test/empty-clone-cd.txt: artifacts/build/debug/$(GOHOSTOS)/$(GOHOSTARCH)/grit artifacts/test/normal-clone.txt
	-@mkdir -p "$(@D)"
	GRIT_CONFIG=etc/testing-file.toml "$(<)" cd empty | grep /tmp/grit-test/clone/github.com/git-fixtures/empty
	GRIT_CONFIG=etc/testing-file.toml "$(<)" cd git-fixtures/empty | grep /tmp/grit-test/clone/github.com/git-fixtures/empty

## Testing normal clone change directory
TEST_TARGETS += artifacts/test/normal-clone-cd.txt
artifacts/test/normal-clone-cd.txt: artifacts/build/debug/$(GOHOSTOS)/$(GOHOSTARCH)/grit artifacts/test/normal-clone.txt
	-@mkdir -p "$(@D)"
	GRIT_CONFIG=etc/testing-file.toml "$(<)" cd testrepo | grep /tmp/grit-test/clone/tmp/grit-repos/testorg/testrepo
	GRIT_CONFIG=etc/testing-file.toml "$(<)" cd tmp/grit-repos/testorg/testrepo | grep /tmp/grit-test/clone/tmp/grit-repos/testorg/testrepo

test:: $(TEST_TARGETS)

.PHONY: clean-test
clean-test::
	$(RM) -r /tmp/grit-test/clone
	$(RM) -r artifacts/test

