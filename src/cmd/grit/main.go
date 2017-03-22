package main

import (
	"fmt"
	"os"
	"path"

	"github.com/jmalloc/grit/src/config"
	"github.com/jmalloc/grit/src/index"
	"github.com/jmalloc/grit/src/pathutil"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	homeDir, _ := pathutil.HomeDir()

	app.Name = "grit"
	app.Usage = "Index your Git clones."
	app.Version = "0.1.0"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "The path to the Grit configuration file.",
			EnvVar: "GRIT_CONFIG",
			Value:  path.Join(homeDir, ".grit", "config.toml"),
		},
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
			Name:     "config",
			Usage:    "Print the entire Grit configuration.",
			Category: "deprecated",
			Action:   withConfig(configShow),
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
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}

func loadConfig(ctx *cli.Context) (config.Config, error) {
	return config.Load(ctx.GlobalString("config"))
}

func withConfig(fn func(config.Config, *cli.Context) error) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		c, err := loadConfig(ctx)
		if err != nil {
			return err
		}

		err = fn(c, ctx)

		if _, ok := err.(usageError); ok {
			_ = cli.ShowCommandHelp(ctx, ctx.Command.Name)
			fmt.Fprintln(ctx.App.Writer, "")
		}

		return err
	}
}

func withConfigAndIndex(fn func(config.Config, *index.Index, *cli.Context) error) cli.ActionFunc {
	return withConfig(func(c config.Config, ctx *cli.Context) error {
		idx, err := index.Open(c.Index.Store)
		if err != nil {
			return err
		}
		defer idx.Close()

		return fn(c, idx, ctx)
	})
}

func autocompleteSlug(c *cli.Context) {
	if err := withConfigAndIndex(indexKeys)(c); err != nil {
		panic(err)
	}
}
