package main

import (
	"errors"
	"fmt"

	"github.com/urfave/cli"
)

func clone(ctx *cli.Context) error {
	repo := ctx.Args().First()
	if repo == "" {
		return usageError("not enough arguments")
	}

	providers, err := loadProviders(ctx)
	if err != nil {
		return err
	}

	for _, p := range providers {
		ok, err := p.Clone(repo)
		if err != nil {
			return err
		}

		if ok {
			fmt.Println(p.ClonePath(repo))
			return nil
		}
	}

	return errors.New("repository not found")
}
