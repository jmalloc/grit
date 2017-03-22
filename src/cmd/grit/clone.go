package main

import (
	"fmt"
	"os"
	"sync"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"

	"github.com/jmalloc/grit/src/config"
	"github.com/jmalloc/grit/src/index"
	"github.com/jmalloc/grit/src/repo"
	"github.com/urfave/cli"
)

func clone(c config.Config, idx *index.Index, ctx *cli.Context) error {
	url, err := getCloneURL(c, ctx)
	if err != nil {
		return err
	}

	dir, err := getCloneDir(c, ctx, url)
	if err != nil {
		return err
	}

	_, err = git.PlainClone(dir, false /* isBare */, &git.CloneOptions{URL: url})

	if err == nil || err == git.ErrRepositoryAlreadyExists {
		fmt.Fprintln(ctx.App.Writer, dir)
		return idx.Add(dir, index.All())
	}

	_ = os.RemoveAll(dir)
	return err
}

func getCloneURL(c config.Config, ctx *cli.Context) (string, error) {
	slugOrURL := ctx.Args().First()
	if slugOrURL == "" {
		return "", notEnoughArguments
	}

	if _, err := transport.NewEndpoint(slugOrURL); err == nil {
		if n := ctx.String("source"); n != "" {
			return "", usageError("can not combine --source with a URL")
		}

		return slugOrURL, nil
	}

	if n := ctx.String("source"); n != "" {
		if u, ok := c.Clone.Sources[n]; ok {
			return repo.ResolveURL(u, slugOrURL), nil
		}

		return "", unknownSource(n)
	}

	if url, ok := probeForURL(c, ctx, slugOrURL); ok {
		return url, nil
	}

	return "", silentFailure
}

func probeForURL(c config.Config, ctx *cli.Context, slug string) (string, bool) {
	var wg sync.WaitGroup
	var m sync.Mutex
	sources := map[string]string{}

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
				sources[n] = url
			}
		}(n, u)
	}

	wg.Wait()

	return chooseByKey(ctx.App.Writer, sources)
}

func getCloneDir(c config.Config, ctx *cli.Context, url string) (string, error) {
	if ctx.Bool("golang") {
		return repo.GetGoCloneDir(url)
	}

	return repo.GetCloneDir(c, url)
}
