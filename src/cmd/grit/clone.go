package main

import (
	"fmt"
	"os"

	"gopkg.in/src-d/go-git.v4/plumbing/transport"

	"github.com/jmalloc/grit/src/grit"
	"github.com/urfave/cli"
)

func clone(c *grit.Config, ctx *cli.Context) error {
	slug := ctx.Args().First()
	if slug == "" {
		return usageError("not enough arguments")
	}

	for _, p := range c.Providers {
		if _, err := p.Open(slug); err == nil {
			fmt.Fprintln(ctx.App.Writer, p.Path(slug))
			return nil
		}

		dir, ok, err := p.Clone(os.Stderr, slug)
		if err != nil {
			return err
		}

		if ok {
			fmt.Fprintln(ctx.App.Writer, dir)
			return c.Index.Add(dir)
		}
	}

	return transport.ErrRepositoryNotFound
}

func find(c *grit.Config, ctx *cli.Context) error {
	slug := ctx.Args().First()
	if slug == "" {
		return usageError("not enough arguments")
	}

	dirs, err := c.Index.Find(slug)
	if err != nil {
		return err
	}

	if len(dirs) == 0 {
		return cli.NewExitError("", 1)
	}

	for _, dir := range dirs {
		fmt.Fprintln(ctx.App.Writer, dir)
	}

	return nil
}
