package commands

import (
	"context"

	"github.com/jmalloc/grit/cmd/grit2/internal/di"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "clone",
		Short: "clone a repository",
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
			return di.WithinCommand(
				cmd,
				func(ctx context.Context) error {
					<-ctx.Done()
					return ctx.Err()
				},
			)
		},
	}

	Root.AddCommand(cmd)
}
