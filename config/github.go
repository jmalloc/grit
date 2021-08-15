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

	// Token is the access token used to authenticate with the GitHub API.
	Token string
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
		if b.Result.Name() != "github" {
			// Require the "api" parameter for custom github sources.
			return missingParameterInSource(
				n.Pos,
				b.Result.Name(),
				"api",
			)
		}

		// Otherwise use the default API URL for the implicit github source.
		b.Result.API = DefaultConfig.Sources["github"].(GitHubSource).API
	}

	return nil
}

func (b *gitHubSourceBuilder) VisitParameter(n ast.Parameter) error {
	switch n.Key {
	case "api":
		return parameterAsURL(n, &b.Result.API)
	case "token":
		return parameterAsString(n, &b.Result.Token)
	default:
		return unrecognizedParameterInSource(n, b.Result.Name())
	}
}
