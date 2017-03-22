package main

import (
	"fmt"

	"github.com/jmalloc/grit/src/grit"
	"github.com/jmalloc/grit/src/index"
	"github.com/jmalloc/grit/src/pathutil"
	"github.com/urfave/cli"
)

func indexFind(c grit.Config, idx *index.Index, ctx *cli.Context) error {
	slug := ctx.Args().First()
	if slug == "" {
		return notEnoughArguments
	}

	dirs, err := idx.Find(slug)
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		fmt.Fprintln(ctx.App.Writer, dir)
	}

	return nil
}

func indexKeys(c grit.Config, idx *index.Index, ctx *cli.Context) error {
	keys, err := idx.List(ctx.Args().First())
	if err != nil {
		return err
	}

	for _, k := range keys {
		fmt.Fprintln(ctx.App.Writer, k)
	}

	return nil
}

func indexShow(c grit.Config, idx *index.Index, ctx *cli.Context) error {
	_, err := idx.WriteTo(ctx.App.Writer)
	return err
}

func indexRebuild(c grit.Config, idx *index.Index, ctx *cli.Context) error {
	dirs := []string{c.Index.Root}

	if gosrc, err := pathutil.GoSrc(); err == nil {
		if _, ok := pathutil.RelChild(c.Index.Root, gosrc); ok {
			dirs = append(dirs, gosrc)
		}
	}

	return idx.Rebuild(dirs, index.Known(c))
}
