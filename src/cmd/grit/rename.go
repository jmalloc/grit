package main

import (
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"

	"github.com/jmalloc/grit/src/grit"
	"github.com/jmalloc/grit/src/grit/index"
	"github.com/urfave/cli"
)

func rename(cfg grit.Config, idx *index.Index, c *cli.Context) error {
	slugOrURL := c.Args().First()
	if slugOrURL == "" {
		return errNotEnoughArguments
	}

	src, err := dirFromArg(c, 1)
	if err != nil {
		return err
	}

	base, err := cloneBaseDirFromCurrent(cfg, c, src)
	if err != nil {
		return err
	}

	rem, ok, err := chooseRemote(cfg, c, src, func(rem *config.RemoteConfig, _ transport.Endpoint) string {
		_, u := transformURL(rem, slugOrURL)
		return " --> " + u
	})

	if err != nil {
		return err
	} else if !ok {
		return errSilentFailure
	}

	ep, u := transformURL(rem, slugOrURL)
	rem.URL = u

	err = updateRemote(src, rem)
	if err != nil {
		return err
	}

	dst := grit.EndpointToDir(base, ep)

	return moveClone(cfg, idx, c, src, dst)
}

func transformURL(rem *config.RemoteConfig, slugOrURL string) (ep transport.Endpoint, u string) {
	existing, err := transport.NewEndpoint(rem.URL)
	if err != nil {
		return
	}

	ep, err = transport.NewEndpoint(slugOrURL)
	if err != nil {
		ep = grit.MergeSlug(existing, slugOrURL)
	}

	if grit.EndpointIsSCP(rem.URL) {
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
