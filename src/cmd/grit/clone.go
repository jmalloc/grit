package main

import (
	"fmt"

	"gopkg.in/src-d/go-git.v4/plumbing/transport"

	"github.com/jmalloc/grit/src/config"
	"github.com/urfave/cli"
)

func clone(c *config.Config, ctx *cli.Context) error {
	slug := ctx.Args().First()
	if slug == "" {
		return usageError("not enough arguments")
	}

	for _, p := range c.Providers {
		if _, err := p.Open(slug); err == nil {
			fmt.Println(p.Path(slug))
			return nil
		}

		dir, ok, err := p.Clone(slug)
		if err != nil {
			return err
		}

		if ok {
			fmt.Println(dir)
			return c.Index.Add(dir)
		}
	}

	return transport.ErrRepositoryNotFound
}

func find(c *config.Config, ctx *cli.Context) error {
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
		fmt.Println(dir)
	}

	return nil
}
