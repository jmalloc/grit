package config

import (
	"net/url"

	"github.com/jmalloc/grit/config/internal/ast"
)

// GitHubSource holds the configuration for a repository source that clones
// repositories from a GitHub server.
type GitHubSource struct {
	// SourceName is the name of the source.
	SourceName string

	// API is the URL to the GitHub API.
	API *url.URL
}

func (s GitHubSource) Name() string {
	return s.SourceName
}

func (s GitHubSource) Visit(v SourceVisitor) error {
	return v.VisitGitHubSource(s)
}

// gitHubSourceBuilder is an implementation of ast.Visitor that builds a GitSource
// from the AST.
type gitHubSourceBuilder struct {
	Result GitHubSource
}

func (b *gitHubSourceBuilder) VisitSource(n ast.Source) error {
	b.Result.SourceName = n.Name()

	if err := n.VisitChildren(b); err != nil {
		return err
	}

	if b.Result.API == nil {
		return missingParameterInSource(
			n.Pos,
			b.Result.Name(),
			"api",
		)
	}

	return nil
}

func (b *gitHubSourceBuilder) VisitParameter(n ast.Parameter) error {
	switch n.Key {
	case "api":
		return parameterAsURL(n, &b.Result.API)
	default:
		return unrecognizedParameterInSource(n, b.Result.Name())
	}
}
