package main

import (
	"fmt"

	"github.com/jmalloc/grit/src/config"
	"github.com/jmalloc/grit/src/index"
	"github.com/jmalloc/grit/src/repo"
	"github.com/urfave/cli"
)

func clone(c config.Config, idx *index.Index, ctx *cli.Context) error {
	slug := ctx.Args().First()
	if slug == "" {
		return usageError("not enough arguments")
	}

	var dir string
	var err error

	if ctx.Bool("go") {
		dir, err = repo.CloneToGoPath(c, slug)
	} else {
		dir, err = repo.CloneToCloneRoot(c, slug)
	}

	if err != nil {
		return err
	}

	fmt.Fprintln(ctx.App.Writer, dir)
	return idx.Add(dir, index.All())
}
