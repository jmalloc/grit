package main

import (
	"github.com/jmalloc/grit"
	"github.com/jmalloc/grit/src/grit/index"
	"github.com/jmalloc/grit/src/grit/pathutil"
	"github.com/urfave/cli"
)

func cd(cfg grit.Config, idx *index.Index, c *cli.Context) error {
	if !c.Args().Present() {
		return errNotEnoughArguments
	}

	dir, ok, err := dirFromSlugArg(cfg, idx, c, 0, pathutil.PreferOther)
	if err != nil {
		return err
	} else if !ok {
		return errSilentFailure
	}

	writeln(c, dir)
	exec(c, "cd", dir)
	return nil
}
