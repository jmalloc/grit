package di

import (
	"github.com/jmalloc/grit/config"
	"github.com/jmalloc/grit/source"
)

func init() {
	Provide(func(cfg config.Config) []source.Source {
		var sources []source.Source

		for _, src := range cfg.Sources {
			sources = append(sources, source.FromConfig(src))
		}

		return sources
	})
}
