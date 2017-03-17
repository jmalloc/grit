package provider

import (
	"errors"
	"fmt"
	"os"
	"path"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
)

// Provider is a named Git service.
type Provider struct {
	Name     string
	Driver   Driver
	BasePath string
}

// ClonePath returns the path that r is cloned into.
func (p *Provider) ClonePath(r string) string {
	return path.Join(p.BasePath, r)
}

// Open attempts to the repository r.
func (p *Provider) Open(r string) (*git.Repository, error) {
	url, err := p.Driver.URL(r)
	if err != nil {
		return nil, err
	}

	d := p.ClonePath(r)
	repo, err := git.PlainOpen(d)
	if err != nil {
		return nil, err
	}

	rem, err := repo.Remote("origin")
	if err != nil {
		return nil, err
	}

	if rem.Config().URL != url {
		return nil, fmt.Errorf("found unexpected repo %s at %s", rem.Config().URL, d)
	}

	return repo, nil
}

// Clone clones the repo named r.
func (p *Provider) Clone(r string) (string, bool, error) {
	d := p.ClonePath(r)
	ok, err := p.CloneInto(d, r)
	return d, ok, err
}

// CloneInto clones the repo named r into the directory d.
func (p *Provider) CloneInto(d, r string) (bool, error) {
	url, err := p.Driver.URL(r)
	if err != nil {
		return false, err
	}

	_, err = os.Stat(d)
	if err == nil {
		return false, errors.New("directory already exists")
	} else if !os.IsNotExist(err) {
		return false, err
	}

	_, err = git.PlainClone(d, false, &git.CloneOptions{
		URL: url,
	})

	if err != nil {
		_ = os.RemoveAll(d)

		if err == transport.ErrRepositoryNotFound {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
