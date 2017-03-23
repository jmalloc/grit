package update

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/google/go-github/github"
)

// ErrNoArchive is an error indicating that a release exist but there's no
// archive available for the current platform.
var ErrNoArchive = fmt.Errorf("no release archive found for %s/%s", runtime.GOOS, runtime.GOARCH)

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
	ticker := time.NewTicker(250 * time.Millisecond)
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
