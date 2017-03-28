package main

import (
	"os"
	"path"

	"github.com/jmalloc/grit/src/grit"
	"github.com/jmalloc/grit/src/grit/index"
	"github.com/urfave/cli"
)

func rm(cfg grit.Config, idx *index.Index, c *cli.Context) error {
	dir, err := dirFromFirstArg(c)
	if err != nil {
		return err
	}

	if !c.Bool("force") && !confirm(c, "Are you sure you want to delete this clone?") {
		return errSilentFailure
	}

	if err := os.RemoveAll(dir); err != nil {
		return err
	}

	if c.NArg() == 0 {
		exec(c, "cd", path.Dir(dir))
	}

	return idx.Remove(dir)
}
