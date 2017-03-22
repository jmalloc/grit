package main

import (
	"fmt"

	"github.com/urfave/cli"
)

var (
	notEnoughArguments = usageError("not enough arguments")
	silentFailure      = cli.NewExitError("", 1)
)

type usageError string

func (e usageError) Error() string {
	return string(e)
}

func unknownSource(n string) error {
	return fmt.Errorf("unknown source: %s", n)
}
