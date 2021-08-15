package source

import (
	"context"

	"github.com/jmalloc/grit/internal/config"
)

// Source is an interface for a source of repositories.
type Source interface {
	Name() string
	Type() string
	Description(ctx context.Context) (string, error)
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
