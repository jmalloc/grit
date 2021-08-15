package config

import (
	"net/url"
	"os"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/jmalloc/grit/internal/config/internal/ast"
)

// ParseFile parses configuration from a file.
func ParseFile(filename string) (Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}

	var root ast.Root
	if err := parser.Parse(filename, f, &root); err != nil {
		if e, ok := err.(participle.Error); ok {
			err = newError(
				e.Position(),
				`%s`,
				e.Message(),
			)
		}

		return Config{}, err
	}

	var b builder

	if err := root.Visit(&b); err != nil {
		return Config{}, err
	}

	// Set the default directory if it has not been set.
	if b.Result.Dir == "" {
		b.Result.Dir = DefaultConfig.Dir
	}

	// Add the default sources if they have not been redefined or disabled.
	for _, s := range DefaultConfig.Sources {
		if _, ok := b.SourceDefs[s.Name()]; !ok {
			if b.Result.Sources == nil {
				b.Result.Sources = map[string]Source{}
			}
			b.Result.Sources[s.Name()] = s
		}
	}

	return b.Result, nil
}

// parser produces an AST from config input.
var parser = participle.MustBuild(&ast.Root{})

// builder is an implementation of ast.Visitor that builds a Config from the
// AST.
type builder struct {
	Result     Config
	SourceDefs map[string]lexer.Position
}

// VisitSource adds a source to the configuration.
func (b *builder) VisitSource(n ast.Source) error {
	name := n.Name()

	// Check if there is already a source with this name defined.
	//
	// We check this before checking if the source is disabled so that we can
	// detect when implicit sources have been disabled.
	if pos, ok := b.SourceDefs[name]; ok {
		return newError(
			n.Pos,
			`a source named "%s" is already defined at %s:%d`,
			name,
			pos.Filename,
			pos.Line,
		)
	}

	// Store the position of the source so we can provide accurate information
	// about two different sources with conflicting names.
	if b.SourceDefs == nil {
		b.SourceDefs = map[string]lexer.Position{}
	}
	b.SourceDefs[name] = n.Pos

	// Don't include disable sources in the configuration at all.
	if n.IsDisabled {
		return nil
	}

	src, err := buildSource(n)
	if err != nil {
		return err
	}

	if b.Result.Sources == nil {
		b.Result.Sources = map[string]Source{}
	}
	b.Result.Sources[name] = src

	return nil
}

func buildSource(n ast.Source) (Source, error) {
	switch n.Type {
	case "git":
		var sb gitSourceBuilder
		return sb.Result, n.Visit(&sb)

	case "github":
		var sb gitHubSourceBuilder
		return sb.Result, n.Visit(&sb)

	default:
		return nil, newError(
			n.Pos,
			`unrecognized source type: %s`,
			n.Type,
		)
	}
}

func (b *builder) VisitParameter(n ast.Parameter) error {
	switch n.Key {
	case "dir":
		return parameterAsString(n, &b.Result.Dir)
	default:
		return unrecognizedParameter(n)
	}
}

// parameterAsString gets a non-empty string value from a parameter.
func parameterAsString(p ast.Parameter, target *string) error {
	v, ok := p.Value.AsString()
	if !ok {
		return invalidParameterType(p, "string")
	}

	if v == "" {
		return invalidParameterValue(
			p,
			"value must not be empty",
		)
	}

	*target = v

	return nil
}

// parameterAsURL parses a URL from a parameter value.
func parameterAsURL(p ast.Parameter, target **url.URL) error {
	var v string
	if err := parameterAsString(p, &v); err != nil {
		return err
	}

	var err error
	*target, err = url.Parse(v)
	if err != nil {
		return invalidParameterValue(
			p,
			`expected a URL: %s`,
			err,
		)
	}

	return nil
}

// parameterAsGitEndpoint parses a URL from a parameter value.
func parameterAsGitEndpoint(p ast.Parameter, target **transport.Endpoint) error {
	var v string
	if err := parameterAsString(p, &v); err != nil {
		return err
	}

	// transport.NewEndpoint() currently requires paths to have a slash in them.
	// We work around this by adding "/workaround" to the end of the URL then
	// stripping it back off the path later.
	//
	// See https://github.com/go-git/go-git/pull/324
	var err error
	*target, err = transport.NewEndpoint(v + "/workaround")
	if err != nil {
		return invalidParameterValue(
			p,
			`expected a git endpoint URL: %s`,
			err,
		)
	}

	(*target).Path = strings.TrimSuffix((*target).Path, "/workaround")

	return nil
}
