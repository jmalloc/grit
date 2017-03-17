package grit

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
)

// Provider describes a Git provider.
type Provider struct {
	Type, URL string
}

// Config is the entire Grit configuration.
type Config struct {
	Clone struct {
		Path string
	}
	Providers map[string]Provider
}

// LoadConfig loads the Grit configuration from a file.
func LoadConfig(p string) (c Config, err error) {
	_, err = toml.DecodeFile(p, &c)
	if err != nil && !os.IsNotExist(err) {
		return
	}

	if c.Clone.Path == "" {
		c.Clone.Path = "git"
	}

	c.Clone.Path, err = expandPath(p, c.Clone.Path)
	if err != nil {
		return
	}

	if _, ok := c.Providers["github"]; !ok {
		if c.Providers == nil {
			c.Providers = map[string]Provider{}
		}

		c.Providers["github"] = Provider{"github", "github.com"}
	}

	return c, nil
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
