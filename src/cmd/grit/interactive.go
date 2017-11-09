package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"

	"github.com/jmalloc/grit/src/grit"
	"github.com/urfave/cli"
)

func confirm(c *cli.Context, msg string) bool {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Fprintf(c.App.Writer, "%s [y/n]: ", msg)

		scanner.Scan()
		input := scanner.Text()
		input = strings.TrimSpace(input)
		input = strings.ToLower(input)

		switch input {
		case "y", "yes":
			return true
		case "n", "no", "":
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

	width := len(strconv.Itoa(size))
	f := fmt.Sprintf("  %%%dd) %%s", width)

	for i, o := range opt {
		write(c, f, i+1, o)
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Fprint(c.App.Writer, ": ")

		scanner.Scan()
		input := scanner.Text()

		switch strings.ToLower(input) {
		case "q", "quit":
			return 0, false
		default:
			i64, _ := strconv.ParseUint(input, 10, 64)
			idx := int(i64)

			if idx >= 1 && idx <= size {
				return idx - 1, true
			}
		}
	}
}

func chooseCloneDir(cfg grit.Config, c *cli.Context, dirs []string) (string, bool) {
	var opts []string

	for _, dir := range dirs {
		opts = append(opts, formatDir(cfg, dir))
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
			info := fn(cfg, ep)
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
