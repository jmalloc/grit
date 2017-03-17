package grit

import (
	"errors"
	"fmt"
	"io"
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

// Path returns the path that slug is cloned into.
func (p *Provider) Path(slug string) string {
	return path.Join(p.BasePath, slug)
}

// Open attempts to the repository by slug.
func (p *Provider) Open(slug string) (*git.Repository, error) {
	url, err := p.Driver.URL(slug)
	if err != nil {
		return nil, err
	}

	d := p.Path(slug)
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

// Clone clones a repo to the standard location.
func (p *Provider) Clone(w io.Writer, slug string) (string, bool, error) {
	dir := p.Path(slug)
	ok, err := p.CloneInto(w, dir, slug)
	return dir, ok, err
}

// CloneInto clones a repo to a specific location.
func (p *Provider) CloneInto(w io.Writer, dir, slug string) (bool, error) {
	url, err := p.Driver.URL(slug)
	if err != nil {
		return false, err
	}

	_, err = os.Stat(dir)
	if err == nil {
		return false, errors.New("directory already exists")
	} else if !os.IsNotExist(err) {
		return false, err
	}

	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL:      url,
		Progress: w,
	})

	if err != nil {
		_ = os.RemoveAll(dir)

		if err == transport.ErrRepositoryNotFound {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
