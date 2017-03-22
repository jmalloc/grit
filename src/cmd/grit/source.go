package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/jmalloc/grit/src/config"
	"github.com/jmalloc/grit/src/repo"
	"github.com/urfave/cli"
)

func sourceProbe(c config.Config, ctx *cli.Context) error {
	slug := ctx.Args().First()
	if slug == "" {
		return notEnoughArguments
	}

	var wg sync.WaitGroup
	var m sync.Mutex

	for n, u := range c.Clone.Sources {
		wg.Add(1)
		go func(n, u string) {
			defer wg.Done()
			url := repo.ResolveURL(u, slug)
			ok, err := repo.Exists(url)

			m.Lock()
			defer m.Unlock()

			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else if ok {
				fmt.Fprintln(ctx.App.Writer, n)
			}
		}(n, u)
	}

	wg.Wait()

	return nil
}

func sourceList(c config.Config, ctx *cli.Context) error {
	for _, n := range c.Clone.Order {
		fmt.Fprintln(ctx.App.Writer, n, c.Clone.Sources[n])
	}
	return nil
}
