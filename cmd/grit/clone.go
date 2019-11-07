package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"

	"github.com/jmalloc/grit"
	"github.com/jmalloc/grit/src/grit/index"
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
	r, err := git.PlainClone(dir, false /* isBare */, opts)

	switch err {
	case git.ErrRepositoryAlreadyExists:
		fmt.Fprintln(os.Stderr, "found existing clone")

	case transport.ErrEmptyRemoteRepository:
		fmt.Fprintln(os.Stderr, "cloned an empty repository")

	case nil:
		err = setupTracking(r, dir)
		if err != nil {
			return err
		}

	default:
		_ = os.RemoveAll(dir)
		return err
	}

	writeln(c, dir)
	exec(c, "cd", dir)

	return idx.Add(dir)
}

func setupTracking(r *git.Repository, dir string) error {
	head, err := r.Head()
	if err != nil {
		return err
	}

	if !head.Name().IsBranch() {
		return nil
	}

	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "\n[branch \"%s\"]\n", head.Name().Short())
	fmt.Fprintf(buf, "\tremote = origin\n")
	fmt.Fprintf(buf, "\tmerge = %s\n", head.Name())

	p := path.Join(dir, ".git", "config")
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = buf.WriteTo(f)
	return err
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
	base, err := cloneBaseDir(cfg, c)
	if err != nil {
		return "", err
	}

	target := c.String("target")

	if target == "" {
		return grit.EndpointToDir(base, ep.Normalized), nil
	}

	if c.Bool("golang") {
		return "", usageError("can not combine --target with --golang")
	}

	return filepath.Abs(target)
}
