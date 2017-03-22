package repo

import (
	"path"
	"strings"

	"github.com/jmalloc/grit/src/config"
	"github.com/jmalloc/grit/src/pathutil"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
)

// ResolveURL returns an actual URL from a URL template and slug.
func ResolveURL(url, slug string) string {
	return strings.Replace(url, "*", slug, 1)
}

// GetCloneDir returns the absolute path for a clone of a repository.
func GetCloneDir(c config.Config, url string) (string, error) {
	endpoint, err := transport.NewEndpoint(url)
	if err != nil {
		return "", err
	}

	ext := path.Ext(endpoint.Path)
	p := strings.TrimSuffix(endpoint.Path, ext)

	return path.Join(c.Clone.Root, endpoint.Host+p), nil
}

// GetGoCloneDir returns the absolute path for a clone of a Go repository.
func GetGoCloneDir(url string) (string, error) {
	base, err := pathutil.GoSrc()
	if err != nil {
		return "", err
	}

	endpoint, err := transport.NewEndpoint(url)
	if err != nil {
		return "", err
	}

	ext := path.Ext(endpoint.Path)
	p := strings.TrimSuffix(endpoint.Path, ext)

	return path.Join(base, endpoint.Host+p), nil
}
