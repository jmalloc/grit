package di

import (
	"github.com/spf13/cobra"
)

// WithinCommand invokes fn with arguments populated by the DI container within
// the context of a Cobra CLI command.
func WithinCommand(cmd *cobra.Command, fn interface{}) error {
	return newContainer(
		cmd.Context(),
	).Invoke(fn)
}

// MustWithinCommand invokes fn with arguments populated by the DI container
// within the context of a Cobra CLI command or panics if unable to do so.
func MustWithinCommand(cmd *cobra.Command, fn interface{}) {
	if err := WithinCommand(cmd, fn); err != nil {
		panic(err)
	}
}
