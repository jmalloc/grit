package main

import (
	"net/url"
	"path"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/jmalloc/grit"
	"github.com/jmalloc/grit/index"
	"github.com/jmalloc/grit/pathutil"
	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli"
)

func browse(cfg grit.Config, idx *index.Index, c *cli.Context) error {
	dir, ok, err := dirFromSlugArg(cfg, idx, c, 0, pathutil.PreferBase)
	if err != nil {
		return err
	} else if !ok {
		return nil
	}

	rem, ok, err := chooseRemote(cfg, c, dir, nil)
	if err != nil {
		return err
	} else if !ok {
		return nil
	}

	ep, _, err := grit.EndpointFromRemote(rem)
	if err != nil {
		return err
	}

	u := url.URL{
		Scheme: "https",
		Host:   ep.Host,
		Path:   strings.TrimSuffix(ep.Path, ".git"),
	}

	if err := injectGitHubTreeViewPath(&u, dir); err != nil {
		return err
	}

	writef(c, "opening %s", u.String())

	return open.Run(u.String())
}

func injectGitHubTreeViewPath(u *url.URL, dir string) error {
	// HACK: Assume anything with github in the host is either GitHub.com or a
	// GitHub Enterprise Server installation.
	if !strings.Contains(u.Host, "github") {
		return nil
	}

	r, err := git.PlainOpen(dir)
	if err != nil {
		return err
	}

	head, err := r.Head()
	if err != nil {
		// Most likely a repository with no commits, but we don't want this to
		// prevent the user from opening the repository in the browser.
		return nil
	}

	if branch, err := r.Branch(head.Name().Short()); err == nil {
		if branch.Remote != "" {
			u.Path = path.Join(u.Path, "tree", branch.Name)
		}
	} else if tag, ok := resolveUniqueTag(r, head); ok {
		u.Path = path.Join(u.Path, "tree", tag.Name().Short())
	}

	return nil
}

// resolveUniqueTag returns the reference to a tag that refers to ref, if
// exactly one exists; otherwise, it returns ref.
func resolveUniqueTag(r *git.Repository, ref *plumbing.Reference) (*plumbing.Reference, bool) {
	tags, err := r.Tags()
	if err != nil {
		return ref, false
	}
	defer tags.Close()

	var refs []*plumbing.Reference

	tags.ForEach(
		func(tagRef *plumbing.Reference) error {
			if tagRef.Hash() == ref.Hash() {
				// Lightweight tag that points to ref.
				refs = append(refs, tagRef)
			} else if tag, err := r.TagObject(tagRef.Hash()); err == nil {
				// Annotated tag that points to ref.
				if tag.Target == ref.Hash() {
					refs = append(refs, tagRef)
				}
			}

			return nil
		},
	)

	if len(refs) == 1 {
		return refs[0], true
	}

	return ref, false
}
