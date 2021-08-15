package config

import (
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/jmalloc/grit/internal/config/internal/ast"
)

// GitSource holds the configuration for a repository source that clones
// repositories from a "vanilla" git server.
type GitSource struct {
	// SourceName is the name of the source.
	SourceName string

	// Endpoint is a template for repository endpoint URLs.
	//
	// The URL path is specified using the "level 1" URL template syntax
	// described by RFC-6570.
	//
	// In practice, the string {repo} is replaced with the repository name.
	// Support for "higher levels" of the RFC-6570 specification may be added in
	// the future.
	//
	// See https://datatracker.ietf.org/doc/html/rfc6570.
	Endpoint *transport.Endpoint
}

func (s GitSource) Name() string {
	return s.SourceName
}

func (s GitSource) Visit(v SourceVisitor) error {
	return v.VisitGitSource(s)
}

// gitSourceBuilder is an implementation of ast.Visitor that builds a GitSource
// from the AST.
type gitSourceBuilder struct {
	Result GitSource
}

func (b *gitSourceBuilder) VisitSource(n ast.Source) error {
	b.Result.SourceName = n.Name()

	if err := n.VisitChildren(b); err != nil {
		return err
	}

	if b.Result.Endpoint == nil {
		return missingParameterInSource(
			n.Pos,
			b.Result.Name(),
			"endpoint",
		)
	}

	return nil
}

func (b *gitSourceBuilder) VisitParameter(n ast.Parameter) error {
	switch n.Key {
	case "endpoint":
		return parameterAsGitEndpoint(n, &b.Result.Endpoint)
	default:
		return unrecognizedParameterInSource(n, b.Result.Name())
	}
}
