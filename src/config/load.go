package config

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/jmalloc/grit/src/index"
	"github.com/jmalloc/grit/src/provider"
)

// Config holds Grit configuration.
type Config struct {
	Index     *index.Index
	Providers []*provider.Provider
}

// Load loads the Grit configuration from a file.
func Load(file string) (c *Config, err error) {
	var s schema

	_, err = toml.DecodeFile(file, &s)
	if err != nil && !os.IsNotExist(err) {
		return
	}

	// clone path ...
	if s.Clone.Path == "" {
		s.Clone.Path = "git"
	}
	s.Clone.Path, err = expandPath(file, s.Clone.Path)
	if err != nil {
		return
	}

	// clone order ...
	if len(s.Clone.Order) == 0 {
		s.Clone.Order = []string{"github"}
	}

	// index path ...
	if s.Index.Path == "" {
		s.Index.Path = ".grit/index.db"
	}
	s.Index.Path, err = expandPath(file, s.Index.Path)
	if err != nil {
		return
	}

	c = &Config{}

	// providers ...
	for _, n := range s.Clone.Order {
		var d provider.Driver
		if p, ok := s.Providers[n]; ok {
			d, err = makeDriver(s, p)
		} else if n == "github" {
			d = &provider.GitHubDriver{}
		} else {
			err = fmt.Errorf("unknown provider in clone order: %s", n)
		}

		if err != nil {
			break
		}

		c.Providers = append(c.Providers, &provider.Provider{
			Name:     n,
			Driver:   d,
			BasePath: path.Join(s.Clone.Path, n),
		})
	}

	c.Index, err = index.Open(s.Index.Path, c.Providers)

	return
}

func makeDriver(s schema, p providerSchema) (provider.Driver, error) {
	switch p.Driver {
	case "github":
		return &provider.GitHubDriver{
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
