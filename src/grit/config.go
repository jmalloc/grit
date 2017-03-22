package grit

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/jmalloc/grit/src/pathutil"
)

// Config holds Grit configuration.
type Config struct {
	Clone struct {
		Root    string                      `toml:"root"`
		Sources map[string]EndpointTemplate `toml:"sources"`
	} `toml:"clone"`
	Index struct {
		Root  string `toml:"root"`
		Store string `toml:"store"`
	} `toml:"index"`
}

// LoadConfig loads the Grit configuration from a file.
func LoadConfig(file string) (c Config, err error) {
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
	if err := resolveWithDefault(&c.Clone.Root, base, "~/grit"); err != nil {
		return err
	}

	// add github to the source list if it's not already present ...
	if _, ok := c.Clone.Sources["github"]; !ok {
		c.Clone.Sources["github"] = "git@github.com:{{ .Slug }}.git"
	}

	// check the source URLs are valid ...
	for _, t := range c.Clone.Sources {
		if err := t.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) normalizeIndex(base string) error {
	if err := resolveWithDefault(&c.Index.Root, base, c.Clone.Root); err != nil {
		return err
	}

	return resolveWithDefault(&c.Index.Store, base, "index.db")
}

func resolveWithDefault(p *string, base, def string) (err error) {
	if *p == "" {
		*p = def
	}

	*p, err = pathutil.ResolveFrom(base, *p)
	return
}
