package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jmalloc/grit/src/config"
	"github.com/jmalloc/grit/src/index"
	"github.com/jmalloc/grit/src/pathutil"
	"github.com/urfave/cli"
)

func indexFindCommand(c config.Config, ctx *cli.Context) error {
	slug := ctx.Args().First()
	if slug == "" {
		return usageError("not enough arguments")
	}

	idx, err := index.Open(c.Index.Store)
	if err != nil {
		return err
	}
	defer idx.Close()

	dirs, err := idx.Find(slug)
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		fmt.Fprintln(ctx.App.Writer, dir)
	}

	return nil
}

func indexListCommand(c config.Config, ctx *cli.Context) error {
	idx, err := index.Open(c.Index.Store)
	if err != nil {
		return err
	}
	defer idx.Close()

	keys, err := idx.List(ctx.Args().First())
	if err != nil {
		return err
	}

	for _, k := range keys {
		fmt.Fprintln(ctx.App.Writer, k)
	}

	return nil
}

func indexSelectCommand(c config.Config, ctx *cli.Context) error {
	slug := ctx.Args().First()
	if slug == "" {
		return usageError("not enough arguments")
	}

	idx, err := index.Open(c.Index.Store)
	if err != nil {
		return err
	}
	defer idx.Close()

	dirs, err := idx.Find(slug)
	if err != nil {
		return err
	}

	if len(dirs) == 0 {
		return cli.NewExitError("", 1)
	} else if len(dirs) == 1 {
		fmt.Fprintln(ctx.App.Writer, dirs[0])
		return nil
	}

	gopath, _ := pathutil.GoPath()

	for index, dir := range dirs {
		if gopath != "" && strings.HasPrefix(dir, gopath) {
			rel, err := filepath.Rel(gopath, dir)
			if err == nil && rel[0] != '.' {
				fmt.Fprintf(os.Stderr, "%3d) [go] %s\n", index+1, rel)
				continue
			}
		}

		rel, err := filepath.Rel(c.Clone.Root, dir)
		if err == nil && rel[0] != '.' {
			fmt.Fprintf(os.Stderr, "%3d) [grit] %s\n", index+1, rel)
			continue
		}

		fmt.Fprintf(os.Stderr, "%3d) %s\n", index+1, dir)
	}

	i := promptBetween(1, len(dirs))
	fmt.Fprintf(ctx.App.Writer, dirs[i-1])
	return nil
}

func promptBetween(min, max int) int {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Fprintf(os.Stderr, "> ")

		scanner.Scan()
		input := scanner.Text()

		i64, _ := strconv.ParseUint(input, 10, 64)
		i := int(i64)

		if i >= min && i <= max {
			return i
		}
	}
}

func indexShowCommand(c config.Config, ctx *cli.Context) error {
	idx, err := index.Open(c.Index.Store)
	if err != nil {
		return err
	}
	defer idx.Close()

	_, err = idx.WriteTo(ctx.App.Writer)
	return err
}

func indexRebuildCommand(c config.Config, ctx *cli.Context) error {
	idx, err := index.Open(c.Index.Store)
	if err != nil {
		return err
	}
	defer idx.Close()

	dirs := []string{c.Index.Root}

	if gopath, err := pathutil.GoPath(); err == nil {
		if _, err := filepath.Rel(c.Index.Root, gopath); err == nil {
			dirs = append(dirs, gopath)
		}
	}

	return idx.Rebuild(dirs, index.Known(c))
}
