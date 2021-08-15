package config

import (
	"net/url"
)

// DefaultFile is the default path to the Grit configuration file.
const DefaultFile = "~/.config/grit.conf"

// DefaultConfig is the configuration used if no configuration file is present.
var DefaultConfig = Config{
	Dir: "~/grit",
	Sources: map[string]Source{
		"github": GitHubSource{
			SourceName: "github",
			API: &url.URL{
				Scheme: "https",
				Host:   "api.github.com",
			},
		},
	},
}

// Config is the root of a Grit configuration.
type Config struct {
	Dir     string
	Sources map[string]Source
}

// Source is an interface for the configuration of various repository sources.
type Source interface {
	Name() string
	Visit(SourceVisitor) error
}

// SourceVisitor is an interface for visiting Source configurations.
type SourceVisitor interface {
	VisitGitSource(GitSource) error
	VisitGitHubSource(GitHubSource) error
}
