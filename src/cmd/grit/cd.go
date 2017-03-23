package main

import (
	"fmt"
	"os"

	"github.com/jmalloc/grit/src/grit"
	"github.com/jmalloc/grit/src/index"
	"github.com/jmalloc/grit/src/pathutil"
	"github.com/urfave/cli"
)

func cd(cfg grit.Config, idx *index.Index, c *cli.Context) error {
	slug := c.Args().First()
	if slug == "" {
		return errNotEnoughArguments
	}

	dirs, err := idx.Find(slug)
	if err != nil {
		return err
	}

	gosrc, _ := pathutil.GoSrc()
	var opts []string

	for _, dir := range dirs {
		if rel, ok := pathutil.RelChild(gosrc, dir); ok && gosrc != "" {
			opts = append(opts, fmt.Sprintf("[go] %s", rel))
		} else if rel, ok := pathutil.RelChild(cfg.Clone.Root, dir); ok {
			opts = append(opts, fmt.Sprintf("[grit] %s", rel))
		} else {
			opts = append(opts, dir)
		}
	}

	if i, ok := choose(os.Stderr, opts); ok {
		write(c, dirs[i])
		return nil
	}

	return errSilentFailure
}
