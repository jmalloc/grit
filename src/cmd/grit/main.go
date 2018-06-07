package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/jmalloc/grit/src/cmd/grit/autocomplete"
	"github.com/jmalloc/grit/src/grit"
	"github.com/jmalloc/grit/src/grit/index"
	"github.com/jmalloc/grit/src/grit/pathutil"
	"github.com/jmalloc/grit/src/grit/update"
	isatty "github.com/mattn/go-isatty"
	"github.com/urfave/cli"
)

// VERSION is the current Grit version.
var VERSION = semver.MustParse("0.6.9")

func main() {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		checkForUpdates()
		defer waitForUpdateCheck()
	}

	app := cli.NewApp()
	homeDir, _ := pathutil.HomeDir()

	app.Name = "grit"
	app.Usage = "Index your Git clones."
	app.EnableBashCompletion = true
	app.Before = execOpen
	app.After = execClose
	app.Version = VERSION.String()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "Load configuration from `FILE`.",
			EnvVar: "GRIT_CONFIG",
			Value:  path.Join(homeDir, ".config", "grit.toml"),
		},
		cli.StringFlag{
			Name:   "with-shell-integration",
			Hidden: true,
		},
	}

	var updatePreReleaseFlag cli.Flag

	if update.IsPreRelease(VERSION) {
		app.Version += " (pre-release)"
		// hide the pre-release flag when the current version is a pre-release,
		// but retain it so passing it is not an error.
		updatePreReleaseFlag = &cli.BoolTFlag{
			Name:   "pre-release",
			Hidden: true,
		}
	} else {
		updatePreReleaseFlag = &cli.BoolFlag{
			Name:  "pre-release",
			Usage: "Include pre-releases when searching for latest version.",
		}
	}

	app.Commands = []cli.Command{
		{
			Name:      "clone",
			Usage:     "Clone a repository into a new directory.",
			ArgsUsage: "<slug | url>",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "source, s",
					Usage: "Clone from `<name>` instead of probing all sources.",
				},
				cli.StringFlag{
					Name:  "target, t",
					Usage: "Clone into `<dir>` instead of the default location.",
				},
				cli.BoolFlag{
					Name:  "golang, g",
					Usage: "Clone into the appropriate $GOPATH sub-directory.",
				},
			},
			Action:       withConfigAndIndex(clone),
			BashComplete: autocomplete.New(autocomplete.Slug),
		},
		{
			Name:         "cd",
			Usage:        "Change the current directory to the location of <slug>.",
			ArgsUsage:    "<slug>",
			Action:       withConfigAndIndex(cd),
			BashComplete: autocomplete.New(autocomplete.Slug),
		},
		{
			Name:   "ls",
			Usage:  "List all clones in the index.",
			Action: withConfigAndIndex(ls),
		},
		{
			Name:      "mv",
			Usage:     "Move a clone into the correct directory.",
			ArgsUsage: "[<path>]",
			Action:    withConfigAndIndex(mv),
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "golang, g",
					Usage: "Move into the appropriate $GOPATH sub-directory.",
				},
			},
		},
		{
			Name:         "rm",
			Usage:        "Remove a clone from the filesystem and the index.",
			ArgsUsage:    "[<path>]",
			Action:       withConfigAndIndex(rm),
			BashComplete: autocomplete.New(autocomplete.Slug),
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "force, f",
					Usage: "Do not prompt for confirmation.",
				},
			},
		},
		{
			Name:      "slug",
			Usage:     "Display the index slug for a clone.",
			ArgsUsage: "[<path>]",
			Action:    withConfigAndIndex(slug),
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "all, a",
					Usage: "Include partial slugs.",
				},
				cli.BoolFlag{
					Name:  "exact, e",
					Usage: "Only display slugs if <path> is the clone directory.",
				},
			},
		},
		{
			Name:         "browse",
			Aliases:      []string{"open"},
			Usage:        "Open the repository's home page with the system's default browser.",
			UsageText:    "This command currently assumes all sources refer to GitHub or GitHub Enterprise servers.",
			ArgsUsage:    "[<slug>]",
			Action:       withConfigAndIndex(browse),
			BashComplete: autocomplete.New(autocomplete.Slug),
		},
		{
			Name:      "set-url",
			Usage:     "Set the URL for a Git remote and move the clone into the correct directory.",
			ArgsUsage: "<slug | url> [<path>]",
			Action:    withConfigAndIndex(setURL),
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "golang, g",
					Usage: "Move into the appropriate $GOPATH sub-directory.",
				},
			},
		},
		{
			Name:  "source",
			Usage: "Manage Git sources.",
			Subcommands: []cli.Command{
				{
					Name:         "probe",
					Usage:        "Discover which sources have a repository.",
					ArgsUsage:    "<slug>",
					Action:       withConfig(sourceProbe),
					BashComplete: autocomplete.New(autocomplete.Slug),
				},
				{
					Name:      "ls",
					Usage:     "List the configured sources.",
					ArgsUsage: "[<slug>]",
					Description: "If <slug> is provided the source URLs are rendered as though <slug> were being cloned.\n" +
						"   Otherwise, the URL templates are rendered as they appear in the configuration file.",
					Action: withConfig(sourceList),
				},
			},
		},
		{
			Name:  "index",
			Usage: "Manage the Grit repository index.",
			Subcommands: []cli.Command{
				{
					Name:      "ls",
					Usage:     "List slugs that begin with a prefix.",
					ArgsUsage: "[<prefix>]",
					Action:    withConfigAndIndex(indexList),
				},
				{
					Name:         "find",
					Usage:        "List clone directories for a specific slug.",
					ArgsUsage:    "<slug>",
					Action:       withConfigAndIndex(indexFind),
					BashComplete: autocomplete.New(autocomplete.Slug),
				},
				{
					Name:        "scan",
					Usage:       "Scan the index paths for clone directories.",
					ArgsUsage:   "[<dirs> ...]",
					Description: "If no arguments are provided, the configured index paths are scanned.",
					Action:      withConfigAndIndex(indexScan),
				},
				{
					Name:   "prune",
					Usage:  "Remove directories that no longer exist.",
					Action: withConfigAndIndex(indexPrune),
				},
				{
					Name:   "clear",
					Usage:  "Delete the entire index.",
					Action: withConfigAndIndex(indexClear),
				},
				{
					Name:   "dump",
					Hidden: true,
					Action: withConfigAndIndex(indexDump),
				},
			},
		},
		{
			Name:   "self-update",
			Usage:  "Update to the latest version of Grit.",
			Action: selfUpdate,
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:  "timeout, t",
					Usage: "Timeout after `<time>` seconds.",
					Value: 60,
				},
				&cli.BoolFlag{
					Name:  "force, f",
					Usage: "Replace the current binary even if it's newer than the latest published release.",
				},
				updatePreReleaseFlag,
			},
		},
		{
			Name:   "config",
			Hidden: true,
			Action: withConfig(configShow),
		},
		{
			Name:   "shell-integration",
			Hidden: true,
			Action: shellIntegration,
		},
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}

// withConfig creates a CLI action function that calls fn with the Grit
// config parameter.
func withConfig(fn func(grit.Config, *cli.Context) error) cli.ActionFunc {
	return func(c *cli.Context) error {
		cfg, err := grit.LoadConfig(c.GlobalString("config"))
		if err != nil {
			return err
		}

		err = fn(cfg, c)

		if _, ok := err.(usageError); ok {
			writef(c, "Incorrect Usage: %s\n", err)
			_ = cli.ShowCommandHelp(c, c.Command.Name)
			return errSilentFailure
		}

		return err
	}
}

// withConfigAndIndex creates a CLI action function that calls fn with the Grit
// config and index parameters.
func withConfigAndIndex(fn func(grit.Config, *index.Index, *cli.Context) error) cli.ActionFunc {
	return withConfig(func(cfg grit.Config, c *cli.Context) error {
		idx, err := index.Open(cfg.Index.Store)
		if err != nil {
			return err
		}
		defer idx.Close()

		return fn(cfg, idx, c)
	})
}

// writef prints to the terminal using the app's output writer.
func writef(c *cli.Context, s string, v ...interface{}) {
	fmt.Fprintf(c.App.Writer, s+"\n", v...)
}

// writef prints to the terminal using the app's output writer.
func writeln(c *cli.Context, s string) {
	fmt.Fprintln(c.App.Writer, s)
}

// cloneBaseDir returns $GOPATH/src if --golang was passed, otherwise it
// returns the configured clone root.
func cloneBaseDir(cfg grit.Config, c *cli.Context) (string, error) {
	if c.Bool("golang") {
		return pathutil.GoSrc()
	}

	return cfg.Clone.Root, nil
}

// cloneBaseDirFromCurrent returns $GOPATH/src if p is already a child of
// $GOPATH/src or if --golang was passed, otherwise it returns the configured
// clone root.
func cloneBaseDirFromCurrent(cfg grit.Config, c *cli.Context, p string) (string, error) {
	gosrc, err := pathutil.GoSrc()

	if c.Bool("golang") {
		return gosrc, err
	}

	if err == nil {
		if _, ok := pathutil.RelChild(gosrc, p); ok {
			return gosrc, err
		}
	}

	return cfg.Clone.Root, nil
}

// dirFromArg returns the n-th arg if it set, or the current working directory.
func dirFromArg(c *cli.Context, n int) (string, error) {
	if c.NArg() <= n {
		return os.Getwd()
	}

	return c.Args()[n], nil
}

// dirFromSlugArgs returns the dir for the slug in the n-th arg, if present, or
// the current working directory.
func dirFromSlugArg(
	cfg grit.Config,
	idx *index.Index,
	c *cli.Context,
	n int,
	cwdDist pathutil.Distance,
) (string, bool, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", false, err
	} else if !c.Args().Present() {
		return cwd, true, nil
	}

	slug := c.Args().First()
	dirs := idx.Find(slug)
	if len(dirs) == 0 {
		return "", false, notIndexed(slug)
	}

	pathutil.SortByDistance(cwd, dirs, cwdDist)
	dir, ok := chooseCloneDir(cfg, c, dirs)

	return dir, ok, nil
}

// formatDir returns dir formatted for display.
func formatDir(cfg grit.Config, dir string) string {
	var tags []string

	cwd, _ := os.Getwd()
	abs, _ := filepath.Abs(dir)
	gosrc, _ := pathutil.GoSrc()

	if rel, ok := pathutil.RelChild(gosrc, dir); ok && gosrc != "" {
		dir = rel
		tags = append(tags, "go")
	}

	if rel, ok := pathutil.RelChild(cfg.Clone.Root, dir); ok {
		dir = rel
		tags = append(tags, "grit")
	}

	if abs == cwd {
		tags = append(tags, "current")
	}

	if len(tags) > 0 {
		return fmt.Sprintf(
			"[%s] %s",
			strings.Join(tags, ", "),
			dir,
		)
	}

	return dir
}
