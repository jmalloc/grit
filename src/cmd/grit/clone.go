package main

import (
	"fmt"

	"github.com/jmalloc/grit/src/clone"
	"github.com/jmalloc/grit/src/config"
	"github.com/jmalloc/grit/src/index"
	"github.com/urfave/cli"
)

func cloneCommand(c config.Config, ctx *cli.Context) error {
	slug := ctx.Args().First()
	if slug == "" {
		return usageError("not enough arguments")
	}

	var dir string
	var err error

	if ctx.Bool("go") {
		dir, err = clone.ToGoPath(c, slug)
	} else {
		dir, err = clone.ToCloneRoot(c, slug)
	}

	if err != nil {
		return err
	}

	idx, err := index.Open(c.Index.Store)
	if err != nil {
		return err
	}
	defer idx.Close()

	fmt.Fprintln(ctx.App.Writer, dir)
	return idx.Add(dir, index.All())
}
