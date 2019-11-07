package main

import (
	"github.com/jmalloc/grit"
	"github.com/jmalloc/grit/src/grit/index"
	"github.com/urfave/cli"
)

func slug(cfg grit.Config, idx *index.Index, c *cli.Context) error {
	dir, err := dirFromArg(c, 0)
	if err != nil {
		return err
	}

	slugs, base := idx.FindByDir(dir)
	if len(slugs) == 0 {
		return errSilentFailure
	}

	if c.Bool("exact") && dir != base {
		return errSilentFailure
	}

	if c.Bool("all") {
		for _, s := range slugs {
			writeln(c, s)
		}
	} else {
		longest := ""
		for _, s := range slugs {
			if len(s) > len(longest) {
				longest = s
			}
		}
		writeln(c, longest)
	}

	return nil
}
