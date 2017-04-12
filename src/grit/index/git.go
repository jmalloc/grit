package index

import (
	"path"

	"github.com/jmalloc/grit/src/grit"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
)

func slugsFromClone(dir string, filter EndpointFilter) (set, error) {
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

	slugs := newSet()
	for _, rem := range remotes {
		ep, err := transport.NewEndpoint(rem.Config().URL)
		if err != nil {
			continue // skip misconfigured remotes
		}

		if filter == nil || filter(ep) {
			p := grit.EndpointToSlug(ep)
			slugs.Add(p, path.Base(p))
		}
	}

	return slugs, nil
}
