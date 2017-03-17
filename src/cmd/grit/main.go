package main

import (
	"fmt"
	"os"

	"github.com/jmalloc/grit/src/config"
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
			Value:  os.Getenv("HOME") + "/.grit/config",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "clone",
			Usage:     "Clone a git repository.",
			ArgsUsage: "<repo>",
			Action:    action(clone),
		},
		{
			Name:  "index",
			Usage: "Manage the repository index.",
			Subcommands: []cli.Command{
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

func loadConfig(ctx *cli.Context) (*config.Config, error) {
	return config.Load(ctx.GlobalString("config"))
}

func action(fn func(*config.Config, *cli.Context) error) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		c, err := loadConfig(ctx)
		if err != nil {
			return err
		}
		defer c.Index.Close()

		err = fn(c, ctx)

		if _, ok := err.(usageError); ok {
			_ = cli.ShowCommandHelp(ctx, ctx.Command.Name)
			fmt.Println("")
		}

		return err
	}
}

type usageError string

func (e usageError) Error() string {
	return string(e)
}
