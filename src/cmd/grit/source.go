package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/jmalloc/grit/src/grit"
	"github.com/urfave/cli"
)

func sourceProbe(cfg grit.Config, c *cli.Context) error {
	slug := c.Args().First()
	if slug == "" {
		return errNotEnoughArguments
	}

	probeSources(cfg, slug, func(n string, ep grit.Endpoint) {
		write(c, "%s", n)
	})

	return nil
}

func sourceList(cfg grit.Config, c *cli.Context) error {
	if c.NArg() > 0 {
		slug := c.Args()[0]
		for n, t := range cfg.Clone.Sources {
			ep, err := t.Resolve(slug)
			if err == nil {
				write(c, "%s %s", n, ep.Actual)
			} else {
				fmt.Fprintf(os.Stderr, "%s %s", n, err)
			}
		}
	} else {
		for n, t := range cfg.Clone.Sources {
			write(c, "%s %s", n, t)
		}
	}

	return nil
}

func probeSources(
	cfg grit.Config,
	slug string,
	fn func(string, grit.Endpoint),
) {
	var wg sync.WaitGroup
	var m sync.Mutex

	for n, t := range cfg.Clone.Sources {
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
