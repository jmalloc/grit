package main

import (
	"fmt"
	"os"

	"github.com/jmalloc/grit/src/config"
	"github.com/jmalloc/grit/src/grit"
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
			Name:      "clone",
			Usage:     "Clone a git repository.",
			ArgsUsage: "<repo>",
			Action:    action(clone),
		},
		// {
		// 	Name:  "where",
		// 	Usage: "Print the path(s) to a repository.",
		// 	Action: func(ctx *cli.Context) error {
		// 		c, err := loadConfig(ctx)
		// 		if err != nil {
		// 			return err
		// 		}
		// 		spew.Dump(c)
		// 		return nil
		// 	},
		// },
		// {
		// 	Name:  "index",
		// 	Usage: "Manage the index.",
		// 	Subcommands: []cli.Command{
		// 		{
		// 			Name:  "rebuild",
		// 			Usage: "Rebuild the entire index.",
		// 			Action: func(ctx *cli.Context) error {
		// 				return nil
		// 			},
		// 		},
		// 	},
		// },
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}

func loadProviders(ctx *cli.Context) (p []*grit.Provider, err error) {
	p, err = config.Load(ctx.GlobalString("config"))
	return
}

func action(fn cli.ActionFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		err := fn(ctx)

		if _, ok := err.(usageError); ok {
			cli.ShowCommandHelp(ctx, ctx.Command.Name)
			fmt.Println("")
		}

		return err
	}
}

type usageError string

func (e usageError) Error() string {
	return string(e)
}
