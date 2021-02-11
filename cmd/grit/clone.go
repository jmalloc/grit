package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"

	"github.com/jmalloc/grit"
	"github.com/jmalloc/grit/index"
	"github.com/urfave/cli"
)

func clone(cfg grit.Config, idx *index.Index, c *cli.Context) error {
	ep, err := getCloneEndpoint(cfg, c)
	if err != nil {
		return err
	}

	dir, err := getCloneDir(cfg, c, ep)
	if err != nil {
		return err
	}

	opts := &git.CloneOptions{
		URL:      ep.Actual,
		Progress: c.App.Writer,
	}
	_, err = git.PlainClone(dir, false /* isBare */, opts)

	switch err {
	case git.ErrRepositoryAlreadyExists:
		fmt.Fprintln(os.Stderr, "found existing clone")

	case transport.ErrEmptyRemoteRepository:
		fmt.Fprintln(os.Stderr, "cloned an empty repository")

	default:
		_ = os.RemoveAll(dir)
		return err

	case nil:
		// fallthrough ...
	}

	writeln(c, dir)
	exec(c, "cd", dir)

	return idx.Add(dir)
}

func getCloneEndpoint(cfg grit.Config, c *cli.Context) (grit.Endpoint, error) {
	slugOrURL := c.Args().First()
	if slugOrURL == "" {
		return grit.Endpoint{}, errNotEnoughArguments
	}

	source := c.String("source")

	normalized, isURL, err := grit.ParseEndpointOrSlug(slugOrURL)
	if err != nil {
		return grit.Endpoint{}, err
	} else if isURL {
		if source != "" {
			return grit.Endpoint{}, usageError("can not combine --source with a URL")
		}

		return grit.Endpoint{
			Actual:     slugOrURL,
			Normalized: normalized,
		}, nil
	}

	if source != "" {
		if t, ok := cfg.Clone.Sources[source]; ok {
			return t.Resolve(slugOrURL)
		}
		return grit.Endpoint{}, unknownSource(source)
	}

	return probeForURL(cfg, c, slugOrURL)
}

func probeForURL(cfg grit.Config, c *cli.Context, slug string) (grit.Endpoint, error) {
	var sources []string
	var endpoints []grit.Endpoint

	probeSources(cfg, slug, func(n string, ep grit.Endpoint) {
		sources = append(sources, n)
		endpoints = append(endpoints, ep)
	})

	if len(sources) == 0 {
		return grit.Endpoint{}, noSource(slug)
	}

	if i, ok := choose(c, sources); ok {
		return endpoints[i], nil
	}

	return grit.Endpoint{}, errSilentFailure
}

func getCloneDir(cfg grit.Config, c *cli.Context, ep grit.Endpoint) (string, error) {
	target := c.String("target")

	if target == "" {
		return grit.EndpointToDir(cfg.Clone.Root, ep.Normalized), nil
	}

	return filepath.Abs(target)
}

func probeSources(
	cfg grit.Config,
	slug string,
	fn func(string, grit.Endpoint),
) {
	var wg sync.WaitGroup
	var m sync.Mutex

	fmt.Fprintf(os.Stderr, "probing %d source(s) for %s\n", len(cfg.Clone.Sources), slug)

	for n, t := range cfg.Clone.Sources {
		wg.Add(1)

		go func(n string, t grit.EndpointTemplate) {
			defer wg.Done()

			ep, err := t.Resolve(slug)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", n, err)
				return
			}

			fmt.Fprintf(os.Stderr, "%s: trying %s\n", n, ep.Actual)

			exists, err := grit.EndpointExists(ep)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", n, err)
				return
			}

			if exists {
				fmt.Fprintf(os.Stderr, "%s: found %s\n", n, ep.Actual)

				m.Lock()
				defer m.Unlock()
				fn(n, ep)
			}
		}(n, t)
	}

	wg.Wait()
}
