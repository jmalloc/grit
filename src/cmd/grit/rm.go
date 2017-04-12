package main

import (
	"fmt"
	"os"
	"path"

	git "gopkg.in/src-d/go-git.v4"

	"github.com/jmalloc/grit/src/grit"
	"github.com/jmalloc/grit/src/grit/index"
	"github.com/urfave/cli"
)

func rm(cfg grit.Config, idx *index.Index, c *cli.Context) error {
	dir, err := dirFromArg(c, 0)
	if err != nil {
		return err
	}

	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return idx.Remove(dir) // just prune it if it's not on disk
	} else if !info.IsDir() {
		return fmt.Errorf("%s exists but it is not a directory", dir)
	}

	if !c.Bool("force") {
		write(c, "%s:", dir)
		write(c, "")

		if idx.Has(dir) {
			write(c, " ✓ is in the index")
		} else {
			write(c, " - is not in the index")
		}

		uncommitted, err := uncommittedModifications(dir)
		if err != nil {
			write(c, " ✗ does not appear to be a git clone")
		} else if uncommitted == 0 {
			write(c, " ✓ has a clean work tree")
		} else {
			write(c, " ✗ has %d uncommitted modification(s)", uncommitted)
		}

		write(c, "")

		if !confirm(c, "Are you sure you want to delete this directory?") {
			return errSilentFailure
		}
	}

	if err := os.RemoveAll(dir); err != nil {
		return err
	}

	if c.NArg() == 0 {
		exec(c, "cd", path.Dir(dir))
	}

	return idx.Remove(dir)
}

func uncommittedModifications(dir string) (int, error) {
	r, err := git.PlainOpen(dir)
	if err != nil {
		return 0, err
	}

	tree, err := r.Worktree()
	if err == git.ErrIsBareRepository {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	status, err := tree.Status()
	if err != nil {
		return 0, err
	}

	count := 0

	for _, fs := range status {
		if fs.Worktree != git.Unmodified || fs.Staging != git.Unmodified {
			count++
		}
	}

	return count, nil
}
