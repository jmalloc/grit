package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jmalloc/grit/src/config"
	"github.com/jmalloc/grit/src/index"
	"github.com/jmalloc/grit/src/pathutil"
	"github.com/urfave/cli"
)

func cdCommand(c config.Config, ctx *cli.Context) error {
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
