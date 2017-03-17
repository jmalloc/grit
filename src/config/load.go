package config

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/jmalloc/grit/src/grit"
)

// Load loads the Grit configuration from a file.
func Load(file string) (p []*grit.Provider, err error) {
	var c config

	_, err = toml.DecodeFile(file, &c)
	if err != nil && !os.IsNotExist(err) {
		return
	}

	// clone path ...
	if c.Clone.Path == "" {
		c.Clone.Path = "git"
	}
	c.Clone.Path, err = expandPath(file, c.Clone.Path)
	if err != nil {
		return
	}

	// clone order ...
	if len(c.Clone.Order) == 0 {
		c.Clone.Order = []string{"github"}
	}

	// providers ...
	for _, n := range c.Clone.Order {
		var dr grit.Driver
		if pc, ok := c.Providers[n]; ok {
			dr, err = makeDriver(c, pc)
		} else if n == "github" {
			dr = &grit.GitHubDriver{}
		} else {
			err = fmt.Errorf("unknown provider in clone order: %s", n)
		}

		if err != nil {
			break
		}

		p = append(p, &grit.Provider{
			Name:     n,
			Driver:   dr,
			BasePath: path.Join(c.Clone.Path, n),
		})
	}

	return
}

func makeDriver(c config, p provider) (grit.Driver, error) {
	switch p.Driver {
	case "github":
		return &grit.GitHubDriver{
			Host: p.Host,
		}, nil
	default:
		return nil, fmt.Errorf("unknown driver: %s", p.Driver)
	}
}

func expandPath(f, p string) (string, error) {
	if path.IsAbs(p) {
		return p, nil
	}

	base := path.Dir(f)
	if !path.IsAbs(base) {
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}

		base = path.Join(wd, base)
	}

	if !strings.HasPrefix(p, "~/") {
		return path.Join(base, p), nil
	}

	p = strings.TrimPrefix(p, "~/")

	home := os.Getenv("HOME")

	if home == "" {
		return "", errors.New("user home directory is unknown")
	}

	return path.Join(home, p), nil
}
