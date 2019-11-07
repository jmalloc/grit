package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/cavaliercoder/grab"
	humanize "github.com/dustin/go-humanize"
	"github.com/google/go-github/github"
	"github.com/jmalloc/grit/src/grit/update"
	"github.com/urfave/cli"
)

func selfUpdate(c *cli.Context) error {
	// cancel the automatic background check early if it's running
	if updateCheckCancel != nil {
		updateCheckCancel()
	}

	if !c.Bool("force") {
		isBrew, err := isBrewBinary()
		if err != nil {
			return err
		} else if isBrew {
			return errors.New(
				"grit was installed via homebrew, use 'brew upgrade grit' to upgrade",
			)
		}
	}

	// setup a deadline first ...
	timeout := time.Duration(c.Int("timeout")) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	writeln(c, "searching for the latest release")

	gh := github.NewClient(nil)
	preRelease := c.Bool("pre-release")
	rel, err := update.FindLatest(ctx, gh, preRelease)
	if err != nil {
		if err == update.ErrReleaseNotFound && !preRelease {
			return errors.New(err.Error() + ", try --pre-release")
		}
		return err
	}

	latest, err := semver.NewVersion(rel.GetTagName())
	if err != nil {
		return err
	}

	cmp := latest.Compare(VERSION)

	if !c.Bool("force") {
		if cmp == 0 {
			writef(c, "current version (%s) is up to date", VERSION)
			return nil
		} else if cmp < 0 {
			return fmt.Errorf(
				"latest published release (%s) is older than the current version (%s)",
				latest,
				VERSION,
			)
		}
	}

	actualBin, err := os.Executable()
	if err != nil {
		return err
	}

	prefix := fmt.Sprintf("downloading version %s", latest)
	message := prefix
	messageLen := len(message)
	fmt.Fprint(c.App.Writer, message)

	archive, err := update.Download(
		ctx,
		grab.DefaultClient,
		rel,
		func(recv, total uint64) {
			r := float64(recv)
			t := float64(total)
			message = fmt.Sprintf(
				"%s (%d%%, %s / %s)",
				prefix,
				int(r/t*100.0),
				humanize.Bytes(recv),
				humanize.Bytes(total),
			)

			fmt.Fprint(c.App.Writer, "\r"+message)

			l := len(message)
			if messageLen > l {
				clr := strings.Repeat(" ", messageLen-l)
				fmt.Fprint(c.App.Writer, clr)
			}
			messageLen = l
		},
	)
	writeln(c, "")
	if err != nil {
		return err
	}

	latestBin := actualBin + "." + latest.String()
	backupBin := actualBin + "." + VERSION.String() + ".backup"

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

	if cmp > 0 {
		writef(c, "upgraded from version %s to %s", VERSION, latest)
	} else if cmp < 0 {
		writef(c, "downgraded from version %s to %s", VERSION, latest)
	} else if VERSION.String() == latest.String() {
		writef(c, "reinstalled version %s", VERSION)
	} else {
		writef(c, "reinstalled version %s as %s", VERSION, latest)
	}

	return os.Remove(backupBin)
}

var (
	updateCheckContext context.Context
	updateCheckCancel  func()
	updateCheckResult  = make(chan *semver.Version, 1)
	updateCheckPeriod  = 24 * time.Hour
)

func checkForUpdates() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	updateCheckContext = ctx
	updateCheckCancel = cancel

	go func() {
		bin, err := os.Executable()
		if err != nil {
			return
		}

		info, err := os.Stat(bin)
		if err != nil {
			return
		}

		if time.Since(info.ModTime()) < updateCheckPeriod {
			updateCheckCancel()
			return
		}

		gh := github.NewClient(nil)
		if latest, ok, _ := update.IsOutdated(updateCheckContext, gh, VERSION); ok {
			updateCheckResult <- latest
		} else {
			close(updateCheckResult)
		}
	}()
}

func waitForUpdateCheck() {
	defer updateCheckCancel()

	select {
	case <-updateCheckContext.Done():
	case version, ok := <-updateCheckResult:
		if ok {
			fmt.Fprintf(
				os.Stderr,
				"\nNOTICE: An update is available, run %s self-update to install version %s.\n",
				os.Args[0],
				version,
			)
		}

		if bin, err := os.Executable(); err == nil {
			now := time.Now()
			_ = os.Chtimes(bin, now, now)
		}
	}
}

func isBrewBinary() (bool, error) {
	bin, err := os.Executable()
	if err != nil {
		return false, err
	}

	bin, err = filepath.EvalSymlinks(bin)
	if err != nil {
		return false, err
	}

	return strings.HasSuffix(bin, "/Cellar/grit/"+VERSION.String()+"/bin/grit"), nil
}
