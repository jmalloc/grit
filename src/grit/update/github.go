package update

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/google/go-github/github"
)

var (
	// ErrReleaseNotFound that no releases could be found.
	ErrReleaseNotFound = errors.New("no releases")

	// ErrNoArchive is an error indicating that a release exist but there's no
	// archive available for the current platform.
	ErrNoArchive = fmt.Errorf("no release archive found for %s/%s", runtime.GOOS, runtime.GOARCH)
)

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

// Download a release archive for the current platform.
func Download(
	ctx context.Context,
	dl *grab.Client,
	rel *github.RepositoryRelease,
	progress func(uint64, uint64),
) (p string, err error) {
	archive, err := asset(rel)
	if err != nil {
		return
	}

	req, err := grab.NewRequest(archive.GetBrowserDownloadURL())
	if err != nil {
		return
	}

	req.RemoveOnError = true
	req.Size = uint64(archive.GetSize())
	req.Filename = downloadPath(rel)

	ready := dl.DoAsync(req)
	var res *grab.Response

	// wait for the download response to become ready, or the context deadline
	select {
	case <-ctx.Done():
		dl.CancelRequest(req)
		err = ctx.Err()
		return
	case res = <-ready:
	}

	// create a ticker for invoking the progress function ...
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	seen := uint64(0)
	for !res.IsComplete() {
		select {
		case <-ctx.Done():
			dl.CancelRequest(req)
			err = ctx.Err()
			return
		case <-ticker.C:
			if res.BytesTransferred() > seen {
				seen = res.BytesTransferred()
				if progress != nil {
					progress(seen, req.Size)
				}
			}
		}
	}

	p = res.Filename
	err = res.Error
	return
}

func asset(rel *github.RepositoryRelease) (*github.ReleaseAsset, error) {
	for _, a := range rel.Assets {
		if a.GetName() == archiveName {
			return &a, nil
		}
	}

	return nil, ErrNoArchive
}

func downloadPath(rel *github.RepositoryRelease) string {
	return path.Join(
		os.TempDir(),
		fmt.Sprintf("grit-%s-%d.update", rel.GetTagName(), rel.GetID()),
	)
}
