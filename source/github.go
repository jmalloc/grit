package source

import (
	"context"
	"net/http"

	"github.com/google/go-github/v38/github"
	"github.com/jmalloc/grit/config"
	"golang.org/x/oauth2"
)

// GitHub is a repository source that accesses repositories on a GitHub server.
type GitHub struct {
	SourceName string
	Client     *github.Client
}

func (s GitHub) Name() string {
	return s.SourceName
}

func (s GitHub) Description() string {
	return "<desc>"
}

func (f *factory) VisitGitHubSource(src config.GitHubSource) error {
	hc := http.DefaultClient
	if src.Token != "" {
		hc = oauth2.NewClient(
			context.Background(),
			oauth2.StaticTokenSource(&oauth2.Token{
				AccessToken: src.Token,
			}),
		)
	}

	gc := github.NewClient(hc)
	gc.BaseURL = src.API

	f.Result = GitHub{
		SourceName: src.Name(),
		Client:     gc,
	}

	return nil
}
