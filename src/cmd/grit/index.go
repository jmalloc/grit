package main

import (
	"github.com/jmalloc/grit/src/grit"
	"github.com/jmalloc/grit/src/index"
	"github.com/urfave/cli"
)

func indexFind(cfg grit.Config, idx *index.Index, c *cli.Context) error {
	slug := c.Args().First()
	if slug == "" {
		return errNotEnoughArguments
	}

	dirs, err := idx.Find(slug)
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		write(c, dir)
	}

	return nil
}

func indexKeys(cfg grit.Config, idx *index.Index, c *cli.Context) error {
	keys, err := idx.List(c.Args().First())
	if err != nil {
		return err
	}

	for _, k := range keys {
		write(c, k)
	}

	return nil
}

func indexShow(cfg grit.Config, idx *index.Index, c *cli.Context) error {
	_, err := idx.WriteTo(c.App.Writer)
	return err
}

func indexRebuild(cfg grit.Config, idx *index.Index, c *cli.Context) error {
	return idx.Rebuild(cfg.Index.Paths, index.Known(cfg))
}
