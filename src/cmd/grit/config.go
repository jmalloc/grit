package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/jmalloc/grit/src/config"
	"github.com/jmalloc/grit/src/pathutil"
	"github.com/urfave/cli"
)

func configShowCommand(c config.Config, ctx *cli.Context) error {
	enc := toml.NewEncoder(ctx.App.Writer)

	if err := enc.Encode(c); err != nil {
		return err
	}

	file, err := pathutil.Resolve(ctx.GlobalString("config"))
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(
		ctx.App.Writer,
		"\nLoaded from %s\n",
		file,
	)

	return err
}
