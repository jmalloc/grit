package source

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/dogmatiq/cosyne"
	"github.com/google/go-github/v38/github"
	"github.com/jmalloc/grit/internal/config"
	"golang.org/x/oauth2"
)

// GitHub is a repository source that accesses repositories on a GitHub server.
type GitHub struct {
	SourceName string
	Client     *github.Client

	once cosyne.Once
	user *github.User
}

func (s *GitHub) Name() string {
	return s.SourceName
}

func (s *GitHub) Type() string {
	return "github"
}

func (s *GitHub) IsDotCom() bool {
	return s.Client.BaseURL.Host == "api.github.com"
}

func (s *GitHub) Description(ctx context.Context) (string, error) {
	if err := s.once.Do(
		ctx,
		func(context.Context) error {
			user, resp, err := s.Client.Users.Get(ctx, "")
			if err != nil {
				if resp.StatusCode == http.StatusUnauthorized {
					return nil
				}

				return err
			}

			s.user = user

			return nil
		},
	); err != nil {
		return "", err
	}

	host := s.Client.BaseURL.Host
	if s.IsDotCom() {
		host = "github.com"
	}

	auth := "unauthenticated"
	if s.user != nil {
		auth = "as " + s.user.GetLogin()
	}

	return fmt.Sprintf(
		"%s (%s)",
		host,
		auth,
	), nil
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

	u, err := url.Parse(src.API.String()) // clone URL
	if err != nil {
		return err
	}

	if !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
	}

	gc := github.NewClient(hc)
	gc.BaseURL = u

	f.Result = &GitHub{
		SourceName: src.Name(),
		Client:     gc,
	}

	return nil
}
