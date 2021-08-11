package di

import (
	"go.uber.org/dig"
)

// container is the singleton dependency injection container.
var container = dig.New()

// Provide registers a new provider with the container.
//
// See container.Provide() for more information.
func Provide(fn interface{}) error {
	return container.Provide(fn)
}

// Invoke invokes fn with arguments supplied by the container.
func Invoke(fn interface{}) error {
	return container.Invoke(fn)
}

// Close closes the container, calling any functions that were deferred via a
// Deferrer.
func Close() error {
	return container.Invoke(func(d *Deferrer) error {
		return d.Close()
	})
}

// provide calls container.Provide(fn) or panics if unable to do so.
func provide(fn interface{}) {
	if err := container.Provide(fn); err != nil {
		panic(err)
	}
}
