package di

import (
	"context"

	"github.com/spf13/cobra"
)

// RunE returns a function that invokes fn with arguments populated by the
// container. The returned function matches the signature of cobra.Command.RunE.
func RunE(fn interface{}) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		Provide(func() (
			context.Context,
			*cobra.Command,
			[]string,
		) {
			return cmd.Context(), cmd, args
		})

		return Invoke(fn)
	}
}
