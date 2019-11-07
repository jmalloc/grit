# Always run tests by default, even if other makefiles are included beforehand.
.DEFAULT_GOAL := test

# Build Go source files from protocol buffers definitions.
GENERATED_FILES += $(PROTO_FILES:.proto=.pb.go)

# GO_SOURCE_FILES is a space separated list of source files that are used by the
# build process.
GO_SOURCE_FILES += $(shell PATH="$(PATH)" git-find '*.go')

# Disable CGO by default.
# See https://golang.org/cmd/cgo
CGO_ENABLED ?= 0

# GO_APP_VERSION is a human-readable string describing the application version.
# If the "main" package has a variable named "version" it is set to this value
# at link time.
GO_APP_VERSION ?= $(GIT_HEAD_COMMITTISH)

# GO_DEBUG_ARGS and GO_RELEASE_ARGS are arguments passed to "go build" for the
# "debug" and "release" targets, respectively.
GO_DEBUG_ARGS   ?= -v -ldflags "-X main.version=$(GO_APP_VERSION)+debug"
GO_RELEASE_ARGS ?= -v -ldflags "-X main.version=$(GO_APP_VERSION) -s -w"

# Build matrix configuration.
#
# GO_MATRIX_OS is a whitespace separated set of operating systems.
# GO_MATRIX_ARCH is a whitespace separated set of CPU architectures.
#
# The build-matrix is constructed from all permutations of GO_MATRIX_OS and
# GO_MATRIX_ARCH. The default is to build only for the OS and architecture
# specified by the GOHOSTOS and GOHOSTARCH environment variables, that is the OS
# and architecture of current system.
GOHOSTOS   := $(shell go env GOHOSTOS)
GOHOSTARCH := $(shell go env GOHOSTARCH)
GO_MATRIX_OS   ?= $(GOHOSTOS)
GO_MATRIX_ARCH ?= $(GOHOSTARCH)

################################################################################

# _GO_COMMAND_PACKAGES is a list of directory names that are expected to contain
# "main" packages. It forms the basis for the executable filenames.
_GO_COMMAND_PACKAGES = $(notdir $(shell find cmd -type d -mindepth 1 -maxdepth 1 2> /dev/null))

# _GO_EXECUTABLES_xxx is a list of executable filenames to produce in a build.
_GO_EXECUTABLES_NIX = $(_GO_COMMAND_PACKAGES)
_GO_EXECUTABLES_WIN = $(addsuffix .exe,$(_GO_COMMAND_PACKAGES))
ifeq ($(GOHOSTOS),windows)
_GO_EXECUTABLES_HOST = $(_GO_EXECUTABLES_WIN)
else
_GO_EXECUTABLES_HOST = $(_GO_EXECUTABLES_NIX)
endif

# _GO_BUILD_PLATFORM_MATRIX_xxx is the cartesian product of all operating
# systems and architectures specified in GO_MATRIX_OS and GO_MATRIX_ARCH.
_GO_BUILD_PLATFORM_MATRIX_ALL  = $(foreach OS,$(GO_MATRIX_OS),$(foreach ARCH,$(GO_MATRIX_ARCH),$(OS)/$(ARCH)))
_GO_BUILD_PLATFORM_MATRIX_NIX  = $(filter-out windows/%,$(_GO_BUILD_PLATFORM_MATRIX_ALL))
_GO_BUILD_PLATFORM_MATRIX_WIN  = $(filter windows/%,$(_GO_BUILD_PLATFORM_MATRIX_ALL))
_GO_BUILD_PLATFORM_MATRIX_HOST = $(GOHOSTOS)/$(GOHOSTARCH)

# _GO_BUILD_MATRIX_xxx is the cartesian product of the platform matrix and the
# executable filenames.
_GO_BUILD_MATRIX_NIX  = $(foreach P,$(_GO_BUILD_PLATFORM_MATRIX_NIX),$(addprefix $(P)/,$(_GO_EXECUTABLES_NIX)))
_GO_BUILD_MATRIX_WIN  = $(foreach P,$(_GO_BUILD_PLATFORM_MATRIX_WIN),$(addprefix $(P)/,$(_GO_EXECUTABLES_WIN)))
_GO_BUILD_MATRIX_HOST = $(foreach P,$(_GO_BUILD_PLATFORM_MATRIX_HOST),$(addprefix $(P)/,$(_GO_EXECUTABLES_HOST)))

# _GO_DEBUG_TARGETS_xxx is the path to the binaries to produce for debug builds.
_GO_DEBUG_TARGETS_ALL    = $(addprefix artifacts/build/debug/,$(_GO_BUILD_MATRIX_NIX) $(_GO_BUILD_MATRIX_WIN))
_GO_DEBUG_TARGETS_HOST   = $(addprefix artifacts/build/debug/,$(_GO_BUILD_MATRIX_HOST))
.SECONDARY: $(_GO_DEBUG_TARGETS_ALL)

# _GO_DEBUG_TARGETS_xxx is the path to the binaries to produce for release builds.
_GO_RELEASE_TARGETS_ALL  = $(addprefix artifacts/build/release/,$(_GO_BUILD_MATRIX_NIX) $(_GO_BUILD_MATRIX_WIN))
_GO_RELEASE_TARGETS_HOST = $(addprefix artifacts/build/release/,$(_GO_BUILD_MATRIX_HOST))
.SECONDARY: $(_GO_RELEASE_TARGETS_HOST)

# Ensure that Linux release binaries are built before attempting to build a Docker image.
DOCKER_BUILD_REQ += $(addprefix artifacts/build/release/linux/amd64/,$(_GO_EXECUTABLES_NIX))

################################################################################

# test --- Executes all go tests in this module.
.PHONY: test
test: $(GENERATED_FILES)
	go test ./...

# coverage --- Produces an HTML coverage report.
.PHONY: coverage
coverage: artifacts/coverage/index.html

# coverage-open --- Opens the HTML coverage report in a browser.
.PHONY: coverage-open
coverage-open: artifacts/coverage/index.html
	open "$<"

# prepare --- Perform tasks that need to be executed before committing. Stacks
# with the "prepare" target form the common makefile.
.PHONY: prepare
prepare:: test
	go fmt ./...
	go mod tidy

# ci --- Builds a machine-readable coverage report. Stacks with the "ci" target
# from the common makefile.
.PHONY: ci
ci:: artifacts/coverage/cover.out

# clean --- Clears the Go test cache. Stacks with the "clean" target from the
# common makefile.
.PHONY: clean
clean::
	go clean -testcache

# build --- Builds debug executable files suitable for execution on this
# machine. It does not require the current OS and architecture to appear in the
# build matrix.
.PHONY: build
build: $(_GO_DEBUG_TARGETS_HOST)

# debug --- Builds debug executable files for all platforms specified in the
# build matrix.
.PHONY: debug
debug: $(_GO_DEBUG_TARGETS_ALL)

# release --- Builds release executable files for all platforms specified in the
# build matrix.
.PHONY: release
release: $(_GO_RELEASE_TARGETS_ALL)

################################################################################

artifacts/coverage/index.html: artifacts/coverage/cover.out
	go tool cover -html="$<" -o "$@"

.PHONY: artifacts/coverage/cover.out # always rebuild
artifacts/coverage/cover.out: $(GENERATED_FILES)
	@mkdir -p $(@D)
	go test -covermode=count -coverprofile=$@ ./...

artifacts/build/%: $(GO_SOURCE_FILES) $(GENERATED_FILES)
	$(eval PARTS := $(subst /, ,$*))
	$(eval BUILD := $(word 1,$(PARTS)))
	$(eval OS    := $(word 2,$(PARTS)))
	$(eval ARCH  := $(word 3,$(PARTS)))
	$(eval BIN   := $(word 4,$(PARTS)))
	$(eval PKG   := $(basename $(BIN)))
	$(eval ARGS  := $(if $(findstring debug,$(BUILD)),$(GO_DEBUG_ARGS),$(GO_RELEASE_ARGS)))

	CGO_ENABLED=$(CGO_ENABLED) GOOS="$(OS)" GOARCH="$(ARCH)" go build $(ARGS) -o "$@" "./cmd/$(PKG)"
