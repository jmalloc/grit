package main

import (
	"fmt"
	"os"
	"path"

	git "github.com/go-git/go-git/v5"

	"github.com/jmalloc/grit"
	"github.com/jmalloc/grit/index"
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
		writef(c, "%s:", dir)
		writeln(c, "")

		if idx.Has(dir) {
			writeln(c, " ✓ is in the index")
		} else {
			writeln(c, " - is not in the index")
		}

		uncommitted, err := uncommittedModifications(dir)
		if err != nil {
			writeln(c, " ✗ does not appear to be a git clone")
		} else if uncommitted == 0 {
			writeln(c, " ✓ has a clean work tree")
		} else {
			writef(c, " ✗ has %d uncommitted modification(s)", uncommitted)
		}

		writeln(c, "")

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
