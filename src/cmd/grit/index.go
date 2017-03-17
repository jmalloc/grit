package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

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

func indexSelect(c *grit.Config, ctx *cli.Context) error {
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
	} else if len(dirs) == 1 {
		for dir := range dirs {
			fmt.Fprintln(ctx.App.Writer, dir)
			break
		}
	}

	var indices []string
	for dir, rel := range dirs {
		indices = append(indices, dir)
		fmt.Fprintf(os.Stderr, "%d. %s\n", len(indices), rel)
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Fprintf(os.Stderr, "> ")

		scanner.Scan()
		input := scanner.Text()

		index64, _ := strconv.ParseUint(input, 10, 8)
		index := int(index64)

		if index > 0 || index <= len(dirs) {
			fmt.Fprintf(ctx.App.Writer, indices[index-1])
			return nil
		}
	}
}

func indexSearch(c *grit.Config, ctx *cli.Context) error {
	slug := ctx.Args().First()
	if slug == "" {
		return usageError("not enough arguments")
	}

	dirs, err := c.Index.Find(slug)
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		fmt.Fprintln(ctx.App.Writer, dir)
	}

	return nil
}

func indexList(c *grit.Config, ctx *cli.Context) error {
	prefix := ctx.Args().First()

	slugs, err := c.Index.List(prefix)
	if err != nil {
		return err
	}

	for _, slug := range slugs {
		fmt.Fprintln(ctx.App.Writer, slug)
	}

	return nil
}
