package di

import (
	"go.uber.org/dig"
)

// container is the singleton dependency injection container.
var container = dig.New()

// Provide registers a new provider with the container.
func Provide(fn interface{}) {
	if err := container.Provide(fn); err != nil {
		panic(err)
	}
}

// Invoke invokes fn with arguments supplied by the container.
func Invoke(fn interface{}) error {
	err := container.Invoke(fn)
	return unwrapError(err)
}

// Close closes the container, calling any functions that were deferred via a
// Deferrer.
func Close() error {
	return container.Invoke(func(d *Deferrer) error {
		return d.Close()
	})
}
