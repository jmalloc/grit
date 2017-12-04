package main

import (
	"net/url"
	"strings"

	"github.com/jmalloc/grit/src/grit"
	"github.com/jmalloc/grit/src/grit/index"
	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli"
)

func browse(cfg grit.Config, idx *index.Index, c *cli.Context) error {
	dir, ok, err := dirFromSlugArg(cfg, idx, c, 0)
	if err != nil {
		return err
	} else if !ok {
		return nil
	}

	rem, ok, err := chooseRemote(cfg, c, dir, nil)
	if err != nil {
		return err
	} else if !ok {
		return nil
	}

	ep, _, err := grit.EndpointFromRemote(rem)
	if err != nil {
		return err
	}

	u := url.URL{
		Scheme: "https",
		Host:   ep.Host(),
		Path:   strings.TrimSuffix(ep.Path(), ".git"),
	}

	write(c, "opening %s", u.String())

	return open.Run(u.String())
}
