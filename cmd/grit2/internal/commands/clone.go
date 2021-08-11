package commands

import (
	"github.com/jmalloc/grit/cmd/grit2/internal/di"
	"github.com/jmalloc/grit/config"
	"github.com/spf13/cobra"
)

// init adds the "clone" command to the root command.
func init() {
	cmd := &cobra.Command{
		Use:   "clone <repo>",
		Short: "clone a remote repository",
		ValidArgsFunction: func(
			cmd *cobra.Command,
			args []string,
			toComplete string,
		) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveDefault
		},
		RunE: func(
			cmd *cobra.Command,
			args []string,
		) error {
			return di.Invoke(func(
				cfg config.Config,
			) error {
				return nil
			})
		},
	}

	root.AddCommand(cmd)
}
