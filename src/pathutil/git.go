package pathutil

import (
	"path"
	"strings"

	"gopkg.in/src-d/go-git.v4/plumbing/transport"
)

// GetClonePath returns the path relative to the grit "clone root" that the
// repository at the given URL should be cloned into.
func GetClonePath(url string) (string, error) {
	endpoint, err := transport.NewEndpoint(url)
	if err != nil {
		return "", err
	}

	ext := path.Ext(endpoint.Path)
	p := strings.TrimSuffix(endpoint.Path, ext)

	return endpoint.Host + p, nil
}

// GetGoPath returns the path (relative to $GOPATH) that the repository at the
// given URL should be cloned into to work properly with the Go build tools.
func GetGoPath(url string) (string, error) {
	return GetClonePath(url)
}
