package main

import (
	"os"

	"github.com/go-git/go-git/v5/plumbing/transport"

	"github.com/jmalloc/grit"
	"github.com/jmalloc/grit/index"
	"github.com/urfave/cli"
)

func indexList(cfg grit.Config, idx *index.Index, c *cli.Context) error {
	for _, s := range idx.ListSlugs(c.Args().First()) {
		writeln(c, s)
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
		writeln(c, dir)
	}

	return nil
}

func indexScan(cfg grit.Config, idx *index.Index, c *cli.Context) error {
	paths := c.Args()

	if c.NArg() == 0 {
		paths = cfg.Index.Paths
	}

	return idx.Scan(
		c.App.Writer,
		func(ep *transport.Endpoint) bool {
			for _, t := range cfg.Clone.Sources {
				if t.IsMatch(ep) {
					return true
				}
			}

			return false
		},
		paths...,
	)
}

func indexPrune(cfg grit.Config, idx *index.Index, c *cli.Context) error {
	return idx.Prune(c.App.Writer)
}

func indexClear(cfg grit.Config, idx *index.Index, c *cli.Context) error {
	return os.Remove(cfg.Index.Store)
}

func indexDump(cfg grit.Config, idx *index.Index, c *cli.Context) error {
	_, err := idx.WriteTo(c.App.Writer)
	return err
}
