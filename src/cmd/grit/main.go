package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/jmalloc/grit/src/cmd/grit/autocomplete"
	"github.com/jmalloc/grit/src/grit"
	"github.com/jmalloc/grit/src/grit/index"
	"github.com/jmalloc/grit/src/grit/pathutil"
	"github.com/jmalloc/grit/src/grit/update"
	"github.com/urfave/cli"
)

// VERSION is the current Grit version.
var VERSION = semver.MustParse("0.4.1")

func main() {
	app := cli.NewApp()
	homeDir, _ := pathutil.HomeDir()

	app.Name = "grit"
	app.Usage = "Index your Git clones."
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "Load configuration from `FILE`.",
			EnvVar: "GRIT_CONFIG",
			Value:  path.Join(homeDir, ".grit", "config.toml"),
		},
		cli.StringFlag{
			Name:   "shell-commands",
			Hidden: true,
		},
	}

	app.Before = func(c *cli.Context) error {
		file := c.String("shell-commands")
		if file == "" {
			return nil
		}

		f, err := os.Create(file)
		if err != nil {
			return err
		}

		app.Metadata["shell-commands"] = f
		return nil
	}

	app.After = func(c *cli.Context) error {
		if f, ok := c.App.Metadata["shell-commands"].(*os.File); ok {
			return f.Close()
		}

		return nil
	}

	app.Version = VERSION.String()
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
			Usage:        "Change the current directory to an indexed clone directory.",
			ArgsUsage:    "<slug>",
			Action:       withConfigAndIndex(cd),
			BashComplete: autocomplete.New(autocomplete.Slug),
		},
		{
			Name:         "rm",
			Usage:        "Remove a clone from the filesystem and the index.",
			ArgsUsage:    "[<slug>]",
			Action:       withConfigAndIndex(rm),
			BashComplete: autocomplete.New(autocomplete.Slug),
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
					Name:      "scan",
					Usage:     "Scan the index paths for clone directories.",
					ArgsUsage: "[<dirs> ...]",
					Action:    withConfigAndIndex(indexScan),
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
					Name:  "force",
					Usage: "Replace the current binary even if it's newer than the latest published release.",
				},
				updatePreReleaseFlag,
			},
		},
		{
			Name:   "config",
			Usage:  "Print the entire Grit configuration.",
			Hidden: true,
			Action: withConfig(configShow),
		},
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}

func withConfig(fn func(grit.Config, *cli.Context) error) cli.ActionFunc {
	return func(c *cli.Context) error {
		cfg, err := grit.LoadConfig(c.GlobalString("config"))
		if err != nil {
			return err
		}

		err = fn(cfg, c)

		if _, ok := err.(usageError); ok {
			write(c, "Incorrect Usage: %s\n", err)
			_ = cli.ShowCommandHelp(c, c.Command.Name)
			return errSilentFailure
		}

		return err
	}
}

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

func write(c *cli.Context, s string, v ...interface{}) {
	fmt.Fprintf(c.App.Writer, s+"\n", v...)
}

func exec(c *cli.Context, v ...string) {
	f, ok := c.App.Metadata["shell-commands"].(*os.File)
	if !ok {
		return
	}

	for _, a := range v {
		a = "'" + strings.Replace(a, "'", `'\''`, -1) + "' "
		if _, err := io.WriteString(f, a); err != nil {
			panic(err)
		}
	}
	if _, err := io.WriteString(f, "\n"); err != nil {
		panic(err)
	}
}
