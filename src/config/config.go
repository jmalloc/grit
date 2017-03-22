package config

import (
	"fmt"
	"os"
	"path"
	"strings"

	"gopkg.in/src-d/go-git.v4/plumbing/transport"

	"github.com/BurntSushi/toml"
	"github.com/jmalloc/grit/src/pathutil"
)

// Config holds Grit configuration.
type Config struct {
	Clone struct {
		Root    string            `toml:"root"`
		Order   []string          `toml:"order"`
		Sources map[string]string `toml:"sources"`
	} `toml:"clone"`
	Index struct {
		Root  string `toml:"root"`
		Store string `toml:"store"`
	} `toml:"index"`
}

// Load loads the Grit configuration from a file.
func Load(file string) (c Config, err error) {
	file, err = pathutil.Resolve(file)
	if err != nil {
		return
	}

	meta, err := toml.DecodeFile(file, &c)
	if err != nil && !os.IsNotExist(err) {
		return
	}

	if keys := meta.Undecoded(); len(keys) != 0 {
		var s []string
		for _, k := range keys {
			s = append(s, k.String())
		}
		err = fmt.Errorf(
			"grit config: unrecognized keys: %s",
			strings.Join(s, ", "),
		)
		return
	}

	dir := path.Dir(file)
	err = c.normalize(dir)

	return
}

func (c *Config) normalize(base string) error {
	if err := c.normalizeClone(base); err != nil {
		return err
	}

	return c.normalizeIndex(base)
}

func (c *Config) normalizeClone(base string) error {
	if err := resolve(&c.Clone.Root, base, "~/grit"); err != nil {
		return err
	}

	// check the source URLs are valid
	var names []string
	for n, u := range c.Clone.Sources {
		if _, err := transport.NewEndpoint(u); err != nil {
			return err
		}

		if n != "github" {
			names = append(names, n)
		}
	}

	// if no clone order is specified, use the defined sources (in any order)
	// and github at the end
	if len(c.Clone.Order) == 0 {
		c.Clone.Order = append(names, "github")
	}

	// ensure that all sources in the clone order actually exist,
	// and automatically create github if not already present
	for _, n := range c.Clone.Order {
		if _, ok := c.Clone.Sources[n]; !ok {
			if n != "github" {
				return fmt.Errorf("grit config: undeclared source '%s' in clone.order", n)
			}
			c.Clone.Sources["github"] = "git@github.com:*.git"
		}
	}

	return nil
}

func (c *Config) normalizeIndex(base string) error {
	if err := resolve(&c.Index.Root, base, c.Clone.Root); err != nil {
		return err
	}

	return resolve(&c.Index.Store, base, "index.db")
}

func resolve(p *string, base, def string) (err error) {
	if *p == "" {
		*p = def
	}

	*p, err = pathutil.ResolveFrom(base, *p)
	return
}
