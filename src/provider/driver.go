package provider

import (
	"fmt"
	"regexp"

	git "gopkg.in/src-d/go-git.v4"
)

// Driver provides information about a specific type of Git provider.
type Driver interface {
	URL(repo string) (string, error)
	Slugs(*git.Repository) ([]string, error)
}

// GitHubDriver is an implementation of Driver for GitHub and GitHub Enterprise.
type GitHubDriver struct {
	Host string
}

// URL gets the URL for a repo slug.
func (d *GitHubDriver) URL(slug string) (string, error) {
	return fmt.Sprintf(gitHubURLFormat, d.host(), slug), nil
}

// Slugs returns the repo "slugs" for a repository.
func (d *GitHubDriver) Slugs(r *git.Repository) (slugs []string, err error) {
	remotes, err := r.Remotes()
	if err != nil {
		return
	}

	for _, rem := range remotes {
		if s, ok := d.slug(rem.Config().URL); ok {
			slugs = append(slugs, s)
		}
	}

	return
}

func (d *GitHubDriver) slug(u string) (string, bool) {
	matches := gitHubURLPattern.FindStringSubmatch(u)
	if len(matches) == 0 || matches[1] != d.host() {
		return "", false
	}

	return matches[2], true
}

func (d *GitHubDriver) host() string {
	if d.Host == "" {
		return "github.com"
	}

	return d.Host
}

var (
	gitHubURLFormat  = "git@%s:%s.git"
	gitHubURLPattern = regexp.MustCompile("^git@(.+?):(.+?).git$")
)
