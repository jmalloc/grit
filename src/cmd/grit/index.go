package main

import (
	"os"

	"github.com/jmalloc/grit/src/grit"
	"github.com/jmalloc/grit/src/grit/index"
	"github.com/urfave/cli"
)

func indexList(cfg grit.Config, idx *index.Index, c *cli.Context) error {
	for _, s := range idx.List(c.Args().First()) {
		write(c, s)
	}

	return nil
}

func indexFind(cfg grit.Config, idx *index.Index, c *cli.Context) error {
	slug := c.Args().First()
	if slug == "" {
		return errNotEnoughArguments
	}

	dirs := idx.Find(slug)

	if len(dirs) == 0 {
		return errSilentFailure
	}

	for _, dir := range dirs {
		write(c, dir)
	}

	return nil
}

func indexScan(cfg grit.Config, idx *index.Index, c *cli.Context) error {
	return idx.Scan(append(
		cfg.Index.Paths,
		c.Args()...,
	)...)
}

func indexPrune(cfg grit.Config, idx *index.Index, c *cli.Context) error {
	return idx.Prune()
}

func indexClear(cfg grit.Config, idx *index.Index, c *cli.Context) error {
	return os.Remove(cfg.Index.Store)
}

func indexDump(cfg grit.Config, idx *index.Index, c *cli.Context) error {
	_, err := idx.WriteTo(c.App.Writer)
	return err
}
