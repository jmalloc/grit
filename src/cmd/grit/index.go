package main

import (
	"os"
	"path"

	"github.com/jmalloc/grit/src/config"
	"github.com/urfave/cli"
)

func indexRebuild(c *config.Config, ctx *cli.Context) error {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		goPath = path.Join(os.Getenv("HOME"), "go")
	}

	return c.Index.Rebuild(goPath)
}

func indexPrint(c *config.Config, ctx *cli.Context) error {
	_, err := c.Index.WriteTo(ctx.App.Writer)
	return err
}
