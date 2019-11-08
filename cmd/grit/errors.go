package main

import (
	"fmt"

	"github.com/urfave/cli"
)

var (
	errNotEnoughArguments = usageError("not enough arguments")
	errSilentFailure      = cli.NewExitError("", 1)
)

type usageError string

func (e usageError) Error() string {
	return string(e)
}

func unknownSource(n string) error {
	return fmt.Errorf("could not find '%s' in the source list", n)
}

func notIndexed(s string) error {
	return fmt.Errorf("could not find '%s' in the index", s)
}

func noSource(s string) error {
	return fmt.Errorf("could not find '%s' at any of the configured sources", s)
}
