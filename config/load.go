package config

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
)

// Load parses the Grit configuration from a file.
func Load(file string) (Config, error) {
	var cfg Config

	meta, err := toml.DecodeFile(file, &cfg)
	if err != nil {
		return Config{}, err
	}

	if keys := meta.Undecoded(); len(keys) != 0 {
		var s []string
		for _, k := range keys {
			s = append(s, k.String())
		}

		return Config{}, fmt.Errorf(
			"unrecognized keys: %s",
			strings.Join(s, ", "),
		)
	}

	if err := normalize(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
