package main

import (
	"errors"
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

	useGoPath := ctx.Bool("go")
	var goPath string
	if useGoPath {
		var ok bool
		goPath, ok = grit.GoPath()
		if !ok {
			return errors.New("could not determine $GOPATH")
		}
	}

	for _, p := range c.Providers {
		var dir string
		var ok bool
		var err error

		if useGoPath {
			dir, ok, err = p.CloneIntoGoPath(os.Stderr, goPath, slug)
		} else {
			dir, ok, err = p.Clone(os.Stderr, slug)
		}

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
