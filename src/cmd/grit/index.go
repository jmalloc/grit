package main

import (
	"fmt"
	"path/filepath"

	"github.com/jmalloc/grit/src/config"
	"github.com/jmalloc/grit/src/index"
	"github.com/jmalloc/grit/src/pathutil"
	"github.com/urfave/cli"
)

func indexFindCommand(c config.Config, ctx *cli.Context) error {
	slug := ctx.Args().First()
	if slug == "" {
		return usageError("not enough arguments")
	}

	idx, err := index.Open(c.Index.Store)
	if err != nil {
		return err
	}
	defer idx.Close()

	dirs, err := idx.Find(slug)
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		fmt.Fprintln(ctx.App.Writer, dir)
	}

	return nil
}

func indexListCommand(c config.Config, ctx *cli.Context) error {
	idx, err := index.Open(c.Index.Store)
	if err != nil {
		return err
	}
	defer idx.Close()

	keys, err := idx.List(ctx.Args().First())
	if err != nil {
		return err
	}

	for _, k := range keys {
		fmt.Fprintln(ctx.App.Writer, k)
	}

	return nil
}

func indexShowCommand(c config.Config, ctx *cli.Context) error {
	idx, err := index.Open(c.Index.Store)
	if err != nil {
		return err
	}
	defer idx.Close()

	_, err = idx.WriteTo(ctx.App.Writer)
	return err
}

func indexRebuildCommand(c config.Config, ctx *cli.Context) error {
	idx, err := index.Open(c.Index.Store)
	if err != nil {
		return err
	}
	defer idx.Close()

	dirs := []string{c.Index.Root}

	if gosrc, err := pathutil.GoSrc(); err == nil {
		if _, err := filepath.Rel(c.Index.Root, gosrc); err == nil {
			dirs = append(dirs, gosrc)
		}
	}

	return idx.Rebuild(dirs, index.Known(c))
}
