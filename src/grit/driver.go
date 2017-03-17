package grit

import (
	"fmt"
	"path"
	"regexp"

	git "gopkg.in/src-d/go-git.v4"
)

// Driver provides information about a specific type of Git provider.
type Driver interface {
	IsValidSlug(string) bool
	URL(slug string) string
	RelativeGoPath(slug string) string
	IndexKeys(*git.Repository) ([]string, error)
}

// GitHubDriver is an implementation of Driver for GitHub and GitHub Enterprise.
type GitHubDriver struct {
	Host string
}

// IsValidSlug returns true if slug is valid for this driver.
func (d *GitHubDriver) IsValidSlug(slug string) bool {
	return gitHubSlugPattern.MatchString(slug)
}

// assertValidSlug panics if slug is invalid.
func (d *GitHubDriver) assertValidSlug(slug string) {
	if !d.IsValidSlug(slug) {
		panic("invalid slug")
	}
}

// URL gets the URL for a repo slug.
func (d *GitHubDriver) URL(slug string) string {
	d.assertValidSlug(slug)
	return fmt.Sprintf(gitHubURLFormat, d.host(), slug)
}

// RelativeGoPath returns the path for a repo relative to $GOPATH.
func (d *GitHubDriver) RelativeGoPath(slug string) string {
	d.assertValidSlug(slug)
	return path.Join("src", d.host(), slug)
}

// IndexKeys returns slugs and any other strings that map to this repo.
func (d *GitHubDriver) IndexKeys(r *git.Repository) (keys []string, err error) {
	remotes, err := r.Remotes()
	if err != nil {
		return
	}

	for _, rem := range remotes {
		keys = append(keys, d.keys(rem.Config().URL)...)
	}

	return
}

func (d *GitHubDriver) keys(u string) []string {
	matches := gitHubURLPattern.FindStringSubmatch(u)

	if len(matches) == 0 || matches[1] != d.host() {
		return nil
	}

	return []string{
		matches[2] + "/" + matches[3],
		matches[3],
	}
}

func (d *GitHubDriver) host() string {
	if d.Host == "" {
		return "github.com"
	}

	return d.Host
}

var (
	gitHubURLFormat   = "git@%s:%s.git"
	gitHubURLPattern  = regexp.MustCompile("^git@(.+?):(.+?)/(.+?).git$")
	gitHubSlugPattern = regexp.MustCompile("^[^/]+/[^/]+$")
)
