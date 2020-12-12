package main

import (
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"

	"github.com/jmalloc/grit"
	"github.com/jmalloc/grit/index"
	"github.com/urfave/cli"
)

func setURL(cfg grit.Config, idx *index.Index, c *cli.Context) error {
	slugOrURL := c.Args().First()
	if slugOrURL == "" {
		return errNotEnoughArguments
	}

	src, err := dirFromArg(c, 1)
	if err != nil {
		return err
	}

	rem, ok, err := chooseRemote(cfg, c, src, func(rem *config.RemoteConfig, _ *transport.Endpoint) string {
		_, u := transformURL(rem, slugOrURL)
		return " --> " + u
	})

	if err != nil {
		return err
	} else if !ok {
		return errSilentFailure
	}

	ep, u := transformURL(rem, slugOrURL)
	rem.URLs = []string{u}

	err = updateRemote(src, rem)
	if err != nil {
		return err
	}

	dst := grit.EndpointToDir(cfg.Clone.Root, ep)

	return moveClone(cfg, idx, c, src, dst)
}

func transformURL(rem *config.RemoteConfig, slugOrURL string) (ep *transport.Endpoint, u string) {
	existing, url, err := grit.EndpointFromRemote(rem)
	if err != nil {
		return
	}

	ep, isURL, err := grit.ParseEndpointOrSlug(slugOrURL)
	if err != nil {
		panic(err)
	} else if !isURL {
		ep = grit.MergeSlug(existing, slugOrURL)
	}

	if grit.EndpointIsSCP(url) {
		u, err = grit.EndpointToSCP(ep)
		if err != nil {
			panic(err)
		}
	} else {
		u = ep.String()
	}

	return
}

func updateRemote(dir string, rem *config.RemoteConfig) error {
	r, err := git.PlainOpen(dir)
	if err != nil {
		return err
	}

	err = r.DeleteRemote(rem.Name)
	if err != nil {
		return err
	}

	_, err = r.CreateRemote(rem)
	return err
}
