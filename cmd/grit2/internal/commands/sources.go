package commands

import (
	"github.com/jmalloc/grit/cmd/grit2/internal/di"
	"github.com/jmalloc/grit/internal/source"
	"github.com/spf13/cobra"
)

// init adds the "clone" command to the root command.
func init() {
	cmd := &cobra.Command{
		Use:   "sources",
		Short: "list information about repository sources",
		Args:  cobra.NoArgs,
		RunE: di.RunE(func(
			cmd *cobra.Command,
			args []string,
			sources []source.Source,
		) error {
			ctx := cmd.Context()

			for _, src := range sources {
				desc, err := src.Description(ctx)
				if err != nil {
					desc = "error: " + err.Error()
				}

				cmd.Printf(
					"%s: %s\n",
					src.Name(),
					desc,
				)
			}

			return nil
		}),
	}

	root.AddCommand(cmd)
}
