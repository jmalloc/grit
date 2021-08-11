// +build !debug

package di

import "go.uber.org/dig"

// unwrapError returns the root cause of err.
func unwrapError(err error) error {
	return dig.RootCause(err)
}
