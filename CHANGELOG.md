# Changelog

## 0.7.4 (2024-10-21)

- **[IMPROVED]** `grit open` (aka `browse`) now opens the GitHub tree view for
  the current branch, tag or hash, instead of always showing the default branch.

## 0.7.3 (2022-03-23)

- The `default-branch` configuration setting now defaults to `main`
- **[FIX]** Fix bug with `grit index scan` that prevented any on-disk clones from being recognized

## 0.7.2 (2021-08-10)

- **[NEW]** Introduce ARM 64 builds for Darwin and Linux, adding Apple Silicon M1 support

## 0.7.1 (2021-03-17)

- **[IMPROVED]** Automatically populate version string during build process

## 0.7.0 (2021-03-16)

- **[FIX]** Fix issue that prevented cloning empty repositories
- **[FIX]** Fix issue where the default branch was added to `.git/config` twice
- **[BC]** Remove `source` sub-command
- **[BC]** Remove `--golang` flag from the `clone` sub-command
- **[BC]** Remove `--golang` flag from the `mv` sub-command
- **[BC]** Remove `--golang` flag from the `set-url` sub-command

## 0.6.10 (2020-12-09)

- **[FIX]** Fix "invalid capabilities" error when cloning from GitHub

## 0.6.9 (2018-06-07)

- **[IMPROVED]** Interactive clone directory prompts now order the options by their distance from the current working directory
- **[IMPROVED]** Interactive clone directory prompts now indicate if one of the options is the current working directory
- **[IMPROVED]** Add `open` as an alias for the `browse` command
- **[FIX]** Naive support for multiple `$GOPATH` directories (grit always used the first one in the list)

## 0.6.8 (2017-12-04)

- **[FIX]** Fix CLI output bug introduced in 0.6.7

## 0.6.7 (2017-12-04)

- **[NEW]** Added `slug` command, which prints the slug for the current clone directory
- **[NEW]** Added `browse` command, which opens the webpage for the current repo (assumes GitHub / GHE)
- **[IMPROVED]** Use `promptui` for interactive CLI elements

## 0.6.6 (2017-11-16)

- **[IMPROVED]** Only check for updates when STDOUT is a TTY
- **[IMPROVED]** Don't allow `self-update` to replace binaries installed by Homebrew

## 0.6.5 (2017-11-15)

- **[NEW]** Added `ls` command, which lists all clones in the index
- **[FIX]** Fix issue that prevented cloning of empty repositories

## 0.6.4 (2017-10-19)

- **[FIX]** Fix issue that prevented Grit from automatically checking for new versions

## 0.6.3 (2017-10-18)

- **[FIX]** Allow proper cloning of symlinks

## 0.6.2 (2017-05-03)

- **[IMPROVED]** Show progress when cloning

## 0.6.1 (2017-04-13)

- **[NEW]** Added `set-url` command, which changes the remote URL then moves the clone accordingly
- **[IMPROVED]** The `mv` command now includes remote information when prompting to choose a directory

## 0.6.0 (2017-04-01)

This release introduces another (and hopefully the final) change to the default
location of the configuration and index files. The configuration file is now
at `~/.config/grit.toml`, following a convention adopted by several other Git
utilities. The index is now stored in the clone root, at `~/grit/index.v2`.
This means that the `~/.grit` directory is no longer used.

If you are using the default locations, you can move your files into the correct
locations by running:

```bash
mkdir -p ~/.config
mv ~/.grit/config.toml ~/.config/grit.toml
mv ~/.grit/index.v2 ~/grit/index.v2
rmdir ~/.grit
```

- **[BC]** The default config location is now `~/.config/grit.toml`
- **[BC]** The default index location is now `~/grit/index.v2`
- **[FIX]** Allow cloning of empty repositories
- **[FIX]** The `index scan` command now accepts relative paths
- **[IMPROVED]** The `rm` command now warns when deleting a clone with uncommitted changes
- **[IMPROVED]** Allow scanning of non-existent paths
- **[IMPROVED]** Print more information when probing sources for repositories

## 0.5.1 (2017-03-28)

- **[FIX]** Fix `self-update` when Grit is invoked via `$PATH`

## 0.5.0 (2017-03-28)

- **[BC]** The `rm` command no longer accepts a slug argument, instead it takes an optional path
- **[NEW]** Added background checks for updates once every 24 hours
- **[NEW]** Added `--force` argument to `rm` to skip confirmation
- **[NEW]** Added `mv` command, which moves an existing clone into the standard clone location
- **[NEW]** Bundled `grit.bash` with the executable, use `eval "$(grit shell-integration)"` in `.bash_profile`
- **[IMPROVED]** When `rm` is invoked with no arguments it changes the current directory to the parent on success

## 0.4.1 (2017-03-27)

- **[NEW]** Added `rm` command, which removes a repository from the filesystem and the index
- **[IMPROVED]** `cd` now prints an error when there are no matching directories
- **[IMPROVED]** `clone` now prints an error when there are no matching sources

## 0.4.0 (2017-03-25)

This release introduces a change to the format of the Grit index. Existing data
in the existing index file will be ignored. The default location for the index
store has also been changed from `~/.grit/index.db` to `~/.grit/index.v2`. If you
are not using the default location for the index store, simply delete the old
data by running `grit index clear` and rebuild the index with `grit index scan`.

- **[BC]** URL template syntax has changed from `{{.Slug}}` to `{{slug}}`
- **[BC]** Removed `index rebuild` command.
- **[BC]** Renamed `index keys` back to `index ls`
- **[BC]** Removed `config` and `index show` commands
- **[FIX]** Usage information is no longer suppressed when running from `grit.bash`
- **[FIX]** Auto-completion of slugs no longer repeats for non-slug parameters
- **[NEW]** Added `index scan` command, which scans the index paths and adds the located repositories to the index
- **[NEW]** Added `index prune` command, which removes non-existent clone directories from the index
- **[NEW]** Added `index clear` command, which erases the entire index
- **[NEW]** Added environment variable substitution in URL templates
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
