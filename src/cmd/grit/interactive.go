package main

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"

	"github.com/jmalloc/grit/src/grit"
	"github.com/jmalloc/grit/src/grit/pathutil"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli"
)

func confirm(c *cli.Context, msg string) bool {
	prompt := promptui.Prompt{
		Label:     msg,
		IsConfirm: true,
	}

	for {
		v, _ := prompt.Run()
		if v == "y" {
			return true
		} else if v == "n" || v == "" {
			return false
		}
	}
}

// choose asks the user to select an entry from opts interactively.
func choose(c *cli.Context, opt []string) (int, bool) {
	size := len(opt)

	if size == 0 {
		return 0, false
	} else if size == 1 {
		return 0, true
	}

	prompt := promptui.Select{
		Label: "",
		Items: opt,
	}

	idx, _, err := prompt.Run()
	return idx, err == nil
}

func chooseCloneDir(cfg grit.Config, c *cli.Context, dirs []string) (string, bool) {
	cwd, _ := os.Getwd()

	// compute "distance from cwd" for each dir
	dists := make([]uint32, len(dirs))
	for idx, dir := range dirs {
		dists[idx] = pathutil.Distance(cwd, dir)
	}

	// sort the dirs such that dirs closest to cwd are listed first
	// any two dirs with the same distance are further sorted by name
	sort.Slice(dirs, func(i, j int) bool {
		di, dj := dists[i], dists[j]

		if di == dj {
			return strings.Compare(dirs[i], dirs[j]) < 0
		}

		return di < dj
	})

	// make options list from the sorted list of dirs
	opts := make([]string, len(dirs))
	for idx, dir := range dirs {
		opts[idx] = formatDir(cfg, dir)
	}

	if i, ok := choose(c, opts); ok {
		return dirs[i], true
	}

	return "", false
}

func chooseRemote(
	cfg grit.Config,
	c *cli.Context,
	dir string,
	fn func(*config.RemoteConfig, transport.Endpoint) string,
) (*config.RemoteConfig, bool, error) {
	r, err := git.PlainOpen(dir)
	if err != nil {
		return nil, false, err
	}

	remotes, err := r.Remotes()
	if err != nil {
		return nil, false, err
	}

	var opts []string
	invalid := map[int]struct{}{}

	for i, rem := range remotes {
		cfg := rem.Config()
		ep, url, err := grit.EndpointFromRemote(cfg)
		if err == nil {
			var info string
			if fn != nil {
				info = fn(cfg, ep)
			}
			opts = append(opts, fmt.Sprintf("[%s] %s%s", cfg.Name, url, info))
		} else {
			opts = append(opts, fmt.Sprintf("[%s] %s (invalid)", cfg.Name, url))
			invalid[i] = struct{}{}
		}
	}

	if i, ok := choose(c, opts); ok {
		if _, ok := invalid[i]; ok {
			return nil, false, errors.New("the selected remote does not have a valid endpoint URL")
		}

		return remotes[i].Config(), true, nil
	}

	return nil, false, nil
}
