package main

import (
	"os"

	"github.com/davecgh/go-spew/spew"
	grit "github.com/jmalloc/grit/src"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "grit"
	app.Usage = "Index your Git clones."
	app.Version = "0.0.0"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "The path to the Grit configuration file.",
			EnvVar: "GRIT_CONFIG",
			Value:  os.Getenv("HOME") + "/.grit.toml",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "where",
			Usage: "Print the path(s) to a repository.",
			Action: func(ctx *cli.Context) error {
				c, err := loadConfig(ctx)
				if err != nil {
					return err
				}
				spew.Dump(c)
				return nil
			},
		},
		{
			Name:  "index",
			Usage: "Manage the index.",
			Subcommands: []cli.Command{
				{
					Name:  "rebuild",
					Usage: "Rebuild the entire index.",
					Action: func(ctx *cli.Context) error {
						return nil
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}

func loadConfig(ctx *cli.Context) (grit.Config, error) {
	return grit.LoadConfig(ctx.GlobalString("config"))
}
