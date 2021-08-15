package source

import (
	"github.com/jmalloc/grit/config"
)

// Source is an interface for a source of repositories.
type Source interface {
	Name() string
	Description() string
}

// FromConfig creates a new source from a source configuration element.
func FromConfig(src config.Source) Source {
	var f factory
	if err := src.Visit(&f); err != nil {
		panic(err)
	}

	return f.Result
}

// factory is an implementation of config.SourceVisitor that constructs sources
// from a config.Source element.
type factory struct {
	Result Source
}
