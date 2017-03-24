package index

import (
	"path"
	"strings"

	"github.com/jmalloc/grit/src/grit"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
)

func slugsFromClone(cfg grit.Config, dir string) (set, error) {
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
		slugs.Merge(slugsFromURL(cfg, rem.Config().URL))
	}

	return slugs, nil
}

func slugsFromURL(cfg grit.Config, url string) set {
	ep, err := transport.NewEndpoint(url)
	if err == nil {
		for _, t := range cfg.Clone.Sources {
			if t.IsMatch(ep) {
				p := strings.TrimSuffix(
					ep.Path[1:],       // trim slash
					path.Ext(ep.Path), // trim .git extension
				)

				return newSet(p, path.Base(p))
			}
		}
	}

	return nil
}

func isGitDir(dir string) bool {
	return isDir(path.Join(dir, ".git"))
}
