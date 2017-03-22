package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/jmalloc/grit/src/grit"
	"github.com/urfave/cli"
)

func sourceProbe(c grit.Config, ctx *cli.Context) error {
	slug := ctx.Args().First()
	if slug == "" {
		return notEnoughArguments
	}

	probeSources(c, slug, func(n string, ep grit.Endpoint) {
		fmt.Fprintln(ctx.App.Writer, n)
	})

	return nil
}

func sourceList(c grit.Config, ctx *cli.Context) error {
	for n, t := range c.Clone.Sources {
		fmt.Fprintln(ctx.App.Writer, n, t)
	}
	return nil
}

func probeSources(
	c grit.Config,
	slug string,
	fn func(string, grit.Endpoint),
) {
	var wg sync.WaitGroup
	var m sync.Mutex

	for n, t := range c.Clone.Sources {
		wg.Add(1)
		go func(n string, t grit.EndpointTemplate) {
			defer wg.Done()

			ep, err := t.Resolve(slug)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}

			exists, err := grit.EndpointExists(ep)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}

			if exists {
				m.Lock()
				defer m.Unlock()
				fn(n, ep)
			}
		}(n, t)
	}

	wg.Wait()
}
