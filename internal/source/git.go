package source

import (
	"context"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/jmalloc/grit/internal/config"
)

// Git is an implementation of the Source interface for "vanilla" Git sources.
type Git struct {
	SourceName string
	Endpoint   *transport.Endpoint
}

func (s *Git) Name() string {
	return s.SourceName
}

func (s *Git) Type() string {
	return "git"
}

func (s *Git) Description(ctx context.Context) (string, error) {
	return s.Endpoint.Host, nil
}

func (f *factory) VisitGitSource(src config.GitSource) error {
	f.Result = &Git{
		SourceName: src.Name(),
		Endpoint:   src.Endpoint,
	}

	return nil
}
