package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"

	"github.com/jmalloc/grit/src/grit"
	"github.com/jmalloc/grit/src/grit/pathutil"
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
	gosrc, _ := pathutil.GoSrc()
	var opts []string

	for _, dir := range dirs {
		if rel, ok := pathutil.RelChild(gosrc, dir); ok && gosrc != "" {
			opts = append(opts, fmt.Sprintf("[go] %s", rel))
		} else if rel, ok := pathutil.RelChild(cfg.Clone.Root, dir); ok {
			opts = append(opts, fmt.Sprintf("[grit] %s", rel))
		} else {
			opts = append(opts, dir)
		}
	}

	if i, ok := choose(c, opts); ok {
		return dirs[i], true
	}

	return "", false
}

func chooseRemote(cfg grit.Config, c *cli.Context, dir string) (*config.RemoteConfig, bool, error) {
	r, err := git.PlainOpen(dir)
	if err != nil {
		return nil, false, err
	}

	remotes, err := r.Remotes()
	if err != nil {
		return nil, false, err
	}

	var opts []string
	for _, rem := range remotes {
		cfg := rem.Config()
		opts = append(opts, fmt.Sprintf("%s %s", cfg.Name, cfg.URL))
	}

	if i, ok := choose(c, opts); ok {
		return remotes[i].Config(), true, nil
	}

	return nil, false, nil
}
