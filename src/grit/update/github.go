package update

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/Masterminds/semver"
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
	rels, _, err := gh.Repositories.ListReleases(ctx, updateRepoOwner, updateRepoName, nil)
	if err != nil {
		return nil, err
	}

	for _, rel := range rels {
		if !preRelease {
			if rel.GetPrerelease() {
				continue
			}

			version, err := semver.NewVersion(rel.GetTagName())
			if err != nil {
				return nil, err
			}

			if IsPreRelease(version) {
				continue
			}
		}

		return rel, nil
	}

	return nil, ErrReleaseNotFound
}

// IsOutdated returns true if v is not the latest available version.
func IsOutdated(
	ctx context.Context,
	gh *github.Client,
	v *semver.Version,
) (latest *semver.Version, outdated bool, err error) {
	rel, err := FindLatest(ctx, gh, IsPreRelease(v))

	if err == nil {
		latest, err = semver.NewVersion(rel.GetTagName())
		outdated = latest.GreaterThan(v)
	}

	return
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

	req, err := grab.NewRequest(
		downloadPath(rel),
		archive.GetBrowserDownloadURL(),
	)
	if err != nil {
		return
	}

	req = req.WithContext(ctx)
	req.Size = int64(archive.GetSize())
	// req.RemoveOnError = true // TODO ??

	res := dl.Do(req)

	// // wait for the download response to become ready, or the context deadline
	// select {
	// case <-ctx.Done():
	// 	dl.CancelRequest(req)
	// 	err = ctx.Err()
	// 	return
	// case res = <-ready:
	// }

	// create a ticker for invoking the progress function ...
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	var seen int64

	for {
		select {
		case <-res.Done:
			p = res.Filename
			err = res.Err()
			return

		case <-ticker.C:
			if res.BytesComplete() > seen {
				seen = res.BytesComplete()
				if progress != nil {
					progress(uint64(seen), uint64(req.Size))
				}
			}
		}
	}
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
