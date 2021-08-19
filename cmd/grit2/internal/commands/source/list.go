package source

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/jmalloc/grit/cmd/grit2/internal/di"
	"github.com/jmalloc/grit/internal/source"
	"github.com/spf13/cobra"
)

// newListCommand returns the "source ls" command.
func newListCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		Short:   "list the configured repository sources",
		Long: heredoc.Doc(`
		The "source list" command lists the configured repository sources and
		their current status.
		`),
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
}
