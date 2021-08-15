package commands

import (
	"errors"

	"github.com/jmalloc/grit/cmd/grit2/internal/di"
	"github.com/jmalloc/grit/source"
	"github.com/spf13/cobra"
)

// init adds the "clone" command to the root command.
func init() {
	cmd := &cobra.Command{
		Use:   "clone <repo>",
		Short: "clone a remote repository",
		Args:  cobra.ExactArgs(1),
		RunE: di.RunE(func(
			cmd *cobra.Command,
			args []string,
			resolver *source.Resolver,
		) error {
			return errors.New("not implemented")
		}),
	}

	root.AddCommand(cmd)
}
