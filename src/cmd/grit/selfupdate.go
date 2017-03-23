package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Masterminds/semver"
	"github.com/cavaliercoder/grab"
	humanize "github.com/dustin/go-humanize"
	"github.com/google/go-github/github"
	"github.com/jmalloc/grit/src/grit/update"
	"github.com/urfave/cli"
)

func selfUpdate(c *cli.Context) error {
	// setup a deadline first ...
	timeout := time.Duration(c.Int("timeout")) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	print(c, "searching for the latest release\n")

	gh := github.NewClient(nil)
	preRelease := c.Bool("pre-release")
	rel, err := update.FindLatest(ctx, gh, preRelease)
	if err != nil {
		if err == update.ErrReleaseNotFound && !preRelease {
			return errors.New(err.Error() + ", try --pre-release")
		}
		return err
	}

	current := semver.MustParse(c.App.Version)
	latest, err := semver.NewVersion(rel.GetTagName())
	if err != nil {
		return err
	}

	if !c.Bool("force") && !latest.GreaterThan(current) {
		return fmt.Errorf(
			"current version (%s) is newer than latest release (%s), not upgrading without --force",
			current,
			latest,
		)
	}

	actualBin, err := filepath.Abs(os.Args[0])
	if err != nil {
		return err
	}

	print(c, "downloading version %s", latest)

	archive, err := update.Download(
		ctx,
		grab.DefaultClient,
		rel,
		func(recv, total uint64) {
			r := float64(recv)
			t := float64(total)

			print(
				c,
				"\rdownloading version %s - %s / %s (%d%%)",
				latest,
				humanize.Bytes(recv),
				humanize.Bytes(total),
				int(r/t*100.0),
			)
		},
	)
	print(c, "\n")
	if err != nil {
		return err
	}

	latestBin := actualBin + "." + latest.String()
	backupBin := actualBin + "." + current.String() + ".backup"

	err = update.Unpack(archive, latestBin)
	if err != nil {
		return err
	}

	err = os.Rename(actualBin, backupBin)
	if err != nil {
		return err
	}

	err = os.Rename(latestBin, actualBin)
	if err != nil {
		return os.Rename(backupBin, actualBin)
	}

	print(c, "updated from v%s to v%s\n", current, latest)
	return os.Remove(backupBin)
}
