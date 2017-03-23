package update

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/go-github/github"
)

// ErrReleaseNotFound that no releases could be found.
var ErrReleaseNotFound = errors.New("no releases")

const (
	updateRepoOwner = "jmalloc"
	updateRepoName  = "grit"
)

// FindLatest finds the latest Grit release.
func FindLatest(ctx context.Context, gh *github.Client, preRelease bool) (*github.RepositoryRelease, error) {
	if !preRelease {
		rel, _, err := gh.Repositories.GetLatestRelease(ctx, updateRepoOwner, updateRepoName)
		if e, ok := err.(*github.ErrorResponse); ok && e.Response.StatusCode == http.StatusNotFound {
			err = ErrReleaseNotFound
		}
		return rel, err
	}

	opts := &github.ListOptions{PerPage: 1}
	rels, _, err := gh.Repositories.ListReleases(ctx, updateRepoOwner, updateRepoName, opts)
	if err != nil {
		return nil, err
	} else if len(rels) == 0 {
		return nil, ErrReleaseNotFound
	}
	return rels[0], nil
}
