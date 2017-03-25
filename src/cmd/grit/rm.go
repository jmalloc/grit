package main

import (
	"os"

	"github.com/jmalloc/grit/src/grit"
	"github.com/jmalloc/grit/src/grit/index"
	"github.com/urfave/cli"
)

func rm(cfg grit.Config, idx *index.Index, c *cli.Context) (err error) {
	dir, err := dirFromArgs(cfg, idx, c)
	if err != nil {
		return err
	}

	write(c, "Are you sure you want to delete this repository?")
	write(c, "Any un-pushed changes will be lost!")
	write(c, "")
	write(c, "\t%s", dir)
	write(c, "")

	if !confirm(c) {
		return errSilentFailure
	}

	if err := os.RemoveAll(dir); err != nil {
		return err
	}

	return idx.Remove(dir)
}

func dirFromArgs(cfg grit.Config, idx *index.Index, c *cli.Context) (string, error) {
	slug := c.Args().First()
	if slug == "" {
		return os.Getwd()
	}

	dirs := idx.Find(slug)
	if len(dirs) == 0 {
		return "", notIndexed(slug)
	}

	if dir, ok := chooseCloneDir(cfg, c, dirs); ok {
		return dir, nil
	}

	return "", errSilentFailure
}
