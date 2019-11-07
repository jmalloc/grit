package main

import (
	"github.com/jmalloc/grit"
	"github.com/jmalloc/grit/src/grit/index"
	"github.com/urfave/cli"
)

func ls(cfg grit.Config, idx *index.Index, c *cli.Context) error {
	dirs := idx.ListClones()

	for _, d := range dirs {
		writeln(c, formatDir(cfg, d))
	}

	return nil
}
