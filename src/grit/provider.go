package grit

import (
	"errors"
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

// StandardPath returns the standard location a repo is cloned to.
func (p *Provider) StandardPath(slug string) string {
	return path.Join(p.BasePath, slug)
}

// RelativeGoPath returns the $GOPATH location for a repo.
func (p *Provider) RelativeGoPath(slug string) string {
	return p.Driver.RelativeGoPath(slug)
}

// Clone clones a repo to the standard location.
func (p *Provider) Clone(w io.Writer, slug string) (string, bool, error) {
	dir := p.StandardPath(slug)
	ok, err := p.CloneInto(w, dir, slug)
	return dir, ok, err
}

// CloneIntoGoPath clones a repo to the appropriate $GOPATH location.
func (p *Provider) CloneIntoGoPath(w io.Writer, goPath, slug string) (string, bool, error) {
	if goPath == "" {
		panic("empty gopath")
	}
	if !p.Driver.IsValidSlug(slug) {
		return "", false, nil
	}

	dir := path.Join(goPath, p.RelativeGoPath(slug))
	ok, err := p.CloneInto(w, dir, slug)
	return dir, ok, err
}

// CloneInto clones a repo to a specific location.
func (p *Provider) CloneInto(w io.Writer, dir, slug string) (bool, error) {
	if !p.Driver.IsValidSlug(slug) {
		return false, nil
	}

	url := p.Driver.URL(slug)

	_, err := os.Stat(dir)
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
