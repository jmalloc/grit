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

	r, err := git.PlainOpen(dir)
	if err != nil {
		return err
	}

	// If we can determine the "HEAD" of the local clone, open the tree view for
	// that commit. Otherwise, open the repository's root.
	if head, err := r.Head(); err == nil {

		// If the head's name is "HEAD", it's a detached HEAD. If the HEAD
		// refers to a branch, the name will be the reference to that branch.
		if head.Name() == "HEAD" {
			// In this case, we check to see if there is a singular tag that
			// refers to the commit, and if so open the tree for the tag name.
			//
			// This is purely for UX, as the user probably expects to see the
			// tag name in the URL.
			head = resolveUniqueTag(r, head)
		}

		// If we still have a detached head, load the tree view for the commit
		// hash; we have no more user-friendly tag or branch name.
		ref := head.Name().Short()
		if ref == "HEAD" {
			ref = head.Hash().String()
		}

		u.Path = path.Join(u.Path, "tree", ref)
	}

	writef(c, "opening %s", u.String())

	return open.Run(u.String())
}

// resolveUniqueTag returns the reference to a tag that refers to ref, if
// exactly one exists; otherwise, it returns ref.
func resolveUniqueTag(r *git.Repository, ref *plumbing.Reference) *plumbing.Reference {
	tags, err := r.Tags()
	if err != nil {
		return ref
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
		return refs[0]
	}

	return ref
}
