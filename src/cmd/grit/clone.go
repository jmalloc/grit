package main

import (
	"fmt"

	"gopkg.in/src-d/go-git.v4/plumbing/transport"

	"github.com/jmalloc/grit/src/config"
	"github.com/urfave/cli"
)

func clone(c *config.Config, ctx *cli.Context) error {
	repo := ctx.Args().First()
	if repo == "" {
		return usageError("not enough arguments")
	}

	for _, p := range c.Providers {
		if _, err := p.Open(repo); err == nil {
			fmt.Println(p.ClonePath(repo))
			return nil
		}

		dir, ok, err := p.Clone(repo)
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
