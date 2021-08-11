package di

import (
	"sync"

	"go.uber.org/multierr"
)

// Deferrer allows DI providers to register functions to be deferred when the
// CLI process ends.
type Deferrer struct {
	m     sync.Mutex
	funcs []func() error
}

// Register a function be executed on exit.
func (d *Deferrer) Defer(fn func() error) {
	d.m.Lock()
	d.funcs = append(d.funcs, fn)
	d.m.Unlock()
}

// Close executes the registered functions in reverse order.
func (d *Deferrer) Close() error {
	d.m.Lock()
	funcs := d.funcs
	d.funcs = nil
	d.m.Unlock()

	var err error

	for _, fn := range funcs {
		fn := fn // capture loop variable

		// Use defer construct so that we get the same panic-handling semantics
		// as usual.
		defer func() {
			err = multierr.Append(
				err,
				fn(),
			)
		}()
	}

	return err
}

func init() {
	Provide(func() *Deferrer {
		return &Deferrer{}
	})
}
