package main

import (
	"fmt"
	"os"
	"path"

	"github.com/jmalloc/grit/src/config"
	"github.com/jmalloc/grit/src/pathutil"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	homeDir, _ := pathutil.HomeDir()

	app.Name = "grit"
	app.Usage = "Index your git clones."
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
			Usage:     "Clone a git repository.",
			ArgsUsage: "<slug>",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "go",
					Usage: "Place the clone under the $GOPATH directory.",
				},
			},
			Action: withConfig(cloneCommand),
		},
		{
			Name:      "cd",
			Usage:     "Interactively prompt for selection of a clone directory.",
			ArgsUsage: "<slug>",
			Action:    withConfig(cdCommand),
			BashComplete: func(c *cli.Context) {
				if err := withConfig(indexListCommand)(c); err != nil {
					panic(err)
				}
			},
		},
		{
			Name:  "config",
			Usage: "Manage the Grit configuration.",
			Subcommands: []cli.Command{
				{
					Name:   "show",
					Usage:  "Display the Grit configuration.",
					Action: withConfig(configShowCommand),
				},
			},
		},
		{
			Name:  "index",
			Usage: "Manage the repository index.",
			Subcommands: []cli.Command{
				{
					Name:      "find",
					Usage:     "List all possible locations for clones of a repository.",
					ArgsUsage: "<slug>",
					Action:    withConfig(indexFindCommand),
					BashComplete: func(c *cli.Context) {
						if err := withConfig(indexListCommand)(c); err != nil {
							panic(err)
						}
					},
				},
				{
					Name:      "list",
					Usage:     "List entries in the index, optionally limited to those matching a prefix.",
					ArgsUsage: "[prefix]",
					Action:    withConfig(indexListCommand),
				},
				{
					Name:   "rebuild",
					Usage:  "Rebuild the entire index.",
					Action: withConfig(indexRebuildCommand),
				},
				{
					Name:   "show",
					Usage:  "Display the entire index database.",
					Action: withConfig(indexShowCommand),
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

type usageError string

func (e usageError) Error() string {
	return string(e)
}
