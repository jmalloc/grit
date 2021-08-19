package source

import (
	"github.com/spf13/cobra"
)

// NewRoot returns the "source" command.
func NewRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "source",
		Short: "manage repository sources",
	}

	cmd.AddCommand(newListCommand())

	return cmd
}
