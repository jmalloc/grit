package commands

import (
	"github.com/jmalloc/grit/cmd/grit2/internal/di"
	"github.com/jmalloc/grit/source"
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
			for _, src := range sources {
				cmd.Printf(
					"%s\t%s\n",
					src.Name(),
					src.Description(),
				)
			}

			return nil
		}),
	}

	root.AddCommand(cmd)
}
