package di

import (
	"context"

	"go.uber.org/dig"
)

// newContainer returns a new DI container.
func newContainer(ctx context.Context, providers ...interface{}) *dig.Container {
	c := dig.New()

	provide(c, func() context.Context {
		return ctx
	})

	for _, p := range providers {
		provide(c, p)
	}

	return c
}

// provide calls c.Provide(p) or panics if unable to do so.
func provide(c *dig.Container, p interface{}) {
	if err := c.Provide(p); err != nil {
		panic(err)
	}
}
