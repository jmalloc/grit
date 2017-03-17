package grit

import (
	"fmt"
)

type Driver interface {
	URL(repo string) (string, error)
}

type GitHubDriver struct {
	Host string
}

func (d *GitHubDriver) URL(repo string) (string, error) {
	return fmt.Sprintf("git@%s:%s.git", d.host(), repo), nil
}

func (d *GitHubDriver) host() string {
	if d.Host == "" {
		return "github.com"
	}

	return d.Host
}
