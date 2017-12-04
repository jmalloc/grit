package main

import (
	"github.com/jmalloc/grit/src/grit"
	"github.com/jmalloc/grit/src/grit/index"
	"github.com/urfave/cli"
)

func cd(cfg grit.Config, idx *index.Index, c *cli.Context) error {
	slug := c.Args().First()
	if slug == "" {
		return errNotEnoughArguments
	}

	dirs := idx.Find(slug)
	if len(dirs) == 0 {
		return notIndexed(slug)
	}

	if dir, ok := chooseCloneDir(cfg, c, dirs); ok {
		writeln(c, dir)
		exec(c, "cd", dir)
		return nil
	}

	return errSilentFailure
}
