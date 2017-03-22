package clone

import (
	"os"
	"path"
	"strings"

	"github.com/jmalloc/grit/src/config"
	"github.com/jmalloc/grit/src/pathutil"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
)

// ToCloneRoot creates a local Git clone of a repository by searching through
// the configured sources.
func ToCloneRoot(c config.Config, slug string) (string, error) {
	return clone(c, slug, c.Clone.Root, pathutil.GetClonePath)
}

// ToGoPath creates a local Git clone of a repository at the appropriate
// location under $GOPATH.
func ToGoPath(c config.Config, slug string) (string, error) {
	p, err := pathutil.GoSrc()
	if err != nil {
		return "", err
	}
	return clone(c, slug, p, pathutil.GetGoPath)
}

func clone(
	c config.Config,
	slug string,
	base string,
	getClonePath func(string) (string, error),
) (string, error) {
	for _, n := range c.Clone.Order {
		url := resolveURL(c.Clone.Sources[n], slug)
		rel, err := getClonePath(url)
		if err != nil {
			continue
		}

		dir := path.Join(base, rel)
		if err := tryClone(url, dir); err == nil {
			return dir, nil
		}
	}

	return "", transport.ErrRepositoryNotFound

}

func resolveURL(url, slug string) string {
	return strings.Replace(url, "*", slug, 1)
}

func tryClone(url, dir string) error {
	if _, err := git.PlainClone(dir, false, &git.CloneOptions{URL: url}); err != nil {
		switch err {
		case git.ErrRepositoryAlreadyExists:
			return nil
		default:
			_ = os.RemoveAll(dir)
			return err
		}
	}

	return nil
}
