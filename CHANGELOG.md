# Changelog

## Next Release

This release introduces a change to the format of the Grit index. The data in
the existing index file will be ignored. The default location for the index
store has also been changed from `~/.grit/index.db` to `~/grit/index.v2`. If you
are not using the default location for the index store, simply delete the old
data by running `grit index clear` and rebuild the index with `grit index scan`.

- **[BC]** Removed `index rebuild` command.
- **[BC]** Renamed `index keys` back to `index ls`
- **[BC]** Removed `config` and `index show` commands
- **[FIX]** Usage information is no longer suppressed when running from `grit.bash`
- **[FIX]** Auto-completion of slugs no longer repeats for non-slug parameters
- **[NEW]** Added `index scan` command, which scans the index paths and adds the located repositories to the index
- **[NEW]** Added `index prune` command, which removes non-existent clone directories from the index
- **[NEW]** Added `index clear` command, which erases the entire index
- **[IMPROVED]** Add slug auto-completion to `clone` command
- **[IMPROVED]** Added the ability to index from arbitrary directories with `index scan`
- **[IMPROVED]** Git submodules are excluded from the index
- **[IMPROVED]** Grit now outputs shell commands to a separate file, see `grit.bash` for details
- **[IMPROVED]** `source ls` command now accepts an optional slug argument for previewing clone URLs

## 0.3.2 (2017-03-24)

- **[BC]** Removed `selfupdate` alias for `self-update` command
- **[FIX]** Backup files are now removed after successful updates

## 0.3.1 (2017-03-24)

- **[IMPROVED]** Better messages about the versions installed by `self-update`

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
- **[BC]** URL templates now uses Go text templates (use `{{.Slug}}` instead of `*`)
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
