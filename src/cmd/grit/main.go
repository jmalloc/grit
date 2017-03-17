package main

import (
	"fmt"
	"os"
	"path"

	"github.com/jmalloc/grit/src/grit"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "grit"
	app.Usage = "Index your git clones."
	app.Version = "0.0.0"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "The path to the Grit configuration file.",
			EnvVar: "GRIT_CONFIG",
			Value:  path.Join(grit.HomeDir(), ".grit", "config.toml"),
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
					Usage: "Place the clone under $GOPATH.",
				},
			},
			Action: action(clone),
		},
		{
			Name:  "index",
			Usage: "Manage the repository index.",
			Subcommands: []cli.Command{
				{
					Name:      "search",
					Usage:     "List the location of all clones of <slug>.",
					ArgsUsage: "<slug>",
					Action:    action(indexSearch),
					BashComplete: func(c *cli.Context) {
						_ = action(indexList)(c)
					},
				},
				{
					Name:      "list",
					Usage:     "List entries in the index.",
					ArgsUsage: "[prefix]",
					Action:    action(indexList),
				},
				{
					Name:   "rebuild",
					Usage:  "Rebuild the index.",
					Action: action(indexRebuild),
				},
				{
					Name:   "print",
					Usage:  "Print the entire index.",
					Action: action(indexPrint),
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}

func loadConfig(ctx *cli.Context) (*grit.Config, error) {
	return grit.LoadConfig(ctx.GlobalString("config"))
}

func action(fn func(*grit.Config, *cli.Context) error) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		c, err := loadConfig(ctx)
		if err != nil {
			return err
		}
		defer c.Index.Close()

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
