# Changelog

## 0.3.0 (2017-03-24)

- **[BC]** Removed `index.root` configuration directive
- **[NEW]** Allow multiple index paths via `index.paths` configuration directive
- **[NEW]** Added `self-update` command
- **[FIX]** Conditionally declare `grit` function in `grit.bash`
- **[FIX]** New clones now setup remote tracking of the default branch

## 0.2.2 (2017-03-23)

- **[FIX]** Fix panic when config file does not exist

## 0.2.1 (2017-03-23)

- **[FIX]** Fix issue where $GOPATH would not be index when outside the index root

## 0.2.0 (2017-03-23)

- **[BC]** Renamed `index select` command to `cd`
- **[BC]** Renamed `clone` command's `--go` flag to `--golang` and added shortcut `-g`
- **[BC]** Removed `clone.order` configuration directive
- **[BC]** Renamed `clone.providers` configuration directive to `clone.sources`
- **[BC]** URL templates now uses Go text templates (use `{{ .Slug }}` instead of `*`)
- **[BC]** Renamed `config show` command to `config` and marked it deprecated
- **[BC]** Rename `index list` command to `index keys`
- **[NEW]** The `clone` command now checks all sources and prompts the user to choose if there are multiple matching repositories
- **[NEW]** Added `--source` flag to `clone`, which specifies a source to use by names
- **[NEW]** Added `--target` flag to `clone`, which specifies target directory for the clone
- **[NEW]** Added `source probe` command, which lists the names of sources that have a given repo
- **[NEW]** Added `source ls` command, which lists the configured sources and their URL templates
- **[NEW]** Added `etc/grit.bash` file, which provides simple shell integration for directory changes and auto-completion
- **[IMPROVED]** `clone` no longer fails if the repository has already been cloned

## 0.1.0 (2017-03-20)

- Initial release
