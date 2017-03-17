package main

import (
	"github.com/jmalloc/grit/src/grit"
	"github.com/urfave/cli"
)

func indexRebuild(c *grit.Config, ctx *cli.Context) error {
	gp, _ := grit.GoPath()
	return c.Index.Rebuild(gp)
}

func indexPrint(c *grit.Config, ctx *cli.Context) error {
	_, err := c.Index.WriteTo(ctx.App.Writer)
	return err
}
