package grit

import (
	"errors"
	"os"
	"path"

	git "gopkg.in/src-d/go-git.v4"
)

type Provider struct {
	Name     string
	Driver   Driver
	BasePath string
}

func (p *Provider) ClonePath(repo string) string {
	return path.Join(p.BasePath, repo)
}

func (p *Provider) Clone(repo string) (bool, error) {
	return p.CloneInto(p.ClonePath(repo), repo)
}

func (p *Provider) CloneInto(dir, repo string) (bool, error) {
	url, err := p.Driver.URL(repo)
	if err != nil {
		return false, err
	}

	if _, err := os.Stat(dir); err != nil {
		if !os.IsNotExist(err) {
			return false, err
		}
	} else {
		return false, errors.New("directory already exists")
	}

	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL: url,
	})

	if err != nil {
		_ = os.RemoveAll(dir)

		if err.Error() == "repository not found" {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
