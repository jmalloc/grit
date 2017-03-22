package index

import (
	"path"
	"strings"

	"github.com/jmalloc/grit/src/config"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
)

// All returns an indexer that indexes any repository.
func All() Indexer {
	return Matching(func(string) bool { return true })
}

// Known returns an indexer that only indexes repositories with URLs known
// to the Grit configuration.
func Known(c config.Config) Indexer {
	return Matching(func(url string) bool {
		a, err := transport.NewEndpoint(url)
		if err != nil {
			return false
		}

		for _, u := range c.Clone.Sources {
			b, _ := transport.NewEndpoint(u)
			if a.Scheme == b.Scheme && a.Host == b.Host {
				return true
			}
		}

		return false
	})
}

// Matching returns an indexer that matches slugs in URLs that match the given
// predicate function.
func Matching(fn func(url string) bool) Indexer {
	return func(dir string) ([]string, error) {
		r, err := git.PlainOpen(dir)
		if err != nil {
			switch err {
			case git.ErrWorktreeNotProvided, git.ErrRepositoryNotExists:
				return nil, nil
			default:
				return nil, err
			}
		}

		remotes, err := r.Remotes()
		if err != nil {
			return nil, err
		}

		var keys []string
		for _, rem := range remotes {
			url := rem.Config().URL
			if fn(url) {
				if k, err := keysFromURL(url); err == nil {
					keys = append(keys, k...)
				}
			}
		}

		return keys, nil
	}
}

// keysFromURL returns a set of keys that map to the repository at the given URL.
func keysFromURL(url string) ([]string, error) {
	endpoint, err := transport.NewEndpoint(url)
	if err != nil {
		return nil, err
	}

	p := strings.TrimPrefix(endpoint.Path, "/")

	ext := path.Ext(p)
	p = strings.TrimSuffix(p, ext)

	return []string{
		p,
		path.Base(p),
	}, nil
}
