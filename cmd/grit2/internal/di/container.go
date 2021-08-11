package di

import (
	"context"

	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

var container = dig.New()
var deferrer Deferrer

func init() {
	provide(func() *Deferrer {
		return &deferrer
	})
}

func Provide(cmd *cobra.Command, fn interface{}) error {
	return container.Provide(fn)
}

func Invoke(cmd *cobra.Command, fn interface{}) error {
	provide(func() context.Context {
		return cmd.Context()
	})

	return container.Invoke(fn)
}

func Close() error {
	return deferrer.Close()
}

// provide calls container.Provide(fn) or panics if unable to do so.
func provide(fn interface{}) {
	if err := container.Provide(fn); err != nil {
		panic(err)
	}
}
