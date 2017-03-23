package main

import (
	"fmt"
	"os"
	"path"

	"github.com/Masterminds/semver"
	"github.com/jmalloc/grit/src/grit"
	"github.com/jmalloc/grit/src/grit/index"
	"github.com/jmalloc/grit/src/grit/pathutil"
	"github.com/jmalloc/grit/src/grit/update"
	"github.com/urfave/cli"
)

// VERSION is the current Grit version.
var VERSION = semver.MustParse("0.3.1")

func main() {
	app := cli.NewApp()
	homeDir, _ := pathutil.HomeDir()

	app.Name = "grit"
	app.Usage = "Index your Git clones."
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "The path to the Grit configuration file.",
			EnvVar: "GRIT_CONFIG",
			Value:  path.Join(homeDir, ".grit", "config.toml"),
		},
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
					Usage: "Clone from a specific named source.",
				},
				cli.StringFlag{
					Name:  "target, t",
					Usage: "Clone into a specific directory.",
				},
				cli.BoolFlag{
					Name:  "golang, g",
					Usage: "Clone into the appropriate $GOPATH sub-directory.",
				},
			},
			Action: withConfigAndIndex(clone),
		},
		{
			Name:         "cd",
			Usage:        "Change the current directory to an indexed clone directory.",
			ArgsUsage:    "<slug>",
			Action:       withConfigAndIndex(cd),
			BashComplete: autocompleteSlug,
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
					BashComplete: autocompleteSlug,
				},
				{
					Name:   "ls",
					Usage:  "List the configured sources.",
					Action: withConfig(sourceList),
				},
			},
		},
		{
			Name:  "index",
			Usage: "Manage the Grit repository index.",
			Subcommands: []cli.Command{
				{
					Name:         "find",
					Usage:        "List directories for a specific repository.",
					ArgsUsage:    "<slug>",
					Action:       withConfigAndIndex(indexFind),
					BashComplete: autocompleteSlug,
				},
				{
					Name:      "keys",
					Usage:     "List keys in the index, optionally limited to those matching a prefix.",
					ArgsUsage: "[prefix]",
					Action:    withConfigAndIndex(indexKeys),
				},
				{
					Name:   "rebuild",
					Usage:  "Rebuild the entire index.",
					Action: withConfigAndIndex(indexRebuild),
				},
				{
					Name:   "show",
					Usage:  "Display the complete repository index.",
					Action: withConfigAndIndex(indexShow),
				},
			},
		},
		{
			Name:    "self-update",
			Aliases: []string{"selfupdate"},
			Usage:   "Update to the latest version of Grit.",
			Action:  selfUpdate,
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:  "timeout, t",
					Usage: "The download timeout, in seconds.",
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
			Name:     "config",
			Usage:    "Print the entire Grit configuration.",
			Category: "deprecated",
			Action:   withConfig(configShow),
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
			_ = cli.ShowCommandHelp(c, c.Command.Name)
			write(c, "")
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

func autocompleteSlug(c *cli.Context) {
	if err := withConfigAndIndex(indexKeys)(c); err != nil {
		panic(err)
	}
}

func write(c *cli.Context, s string, v ...interface{}) {
	fmt.Fprintf(c.App.Writer, s+"\n", v...)
}
