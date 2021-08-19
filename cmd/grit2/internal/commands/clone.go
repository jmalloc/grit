package commands

import (
	"errors"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/jmalloc/grit/cmd/grit2/internal/di"
	"github.com/spf13/cobra"
)

// newCloneCommand returns the "clone" command.
func newCloneCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "clone <repo>",
		Args:  cobra.ExactArgs(1),
		Short: "clone a remote repository",
		Long: heredoc.Doc(`
		The "clone" command makes a local clone of a remote Git repository then
		changes the shell's current working directory to the clone's working
		tree.

		The <repo> argument is a repository name or a part thereof. For example,
		the Grit repository itself may be referred to as "jmalloc/grit" or just
		"grit".

		Each of the repository sources defined in the Grit configuration file is
		searched for matches to the provided repository name. If there are
		multiple matches and the shell is interactive the user is prompted to
		select the desired repository.
		`),
		RunE: di.RunE(func(
			cmd *cobra.Command,
			args []string,
		) error {
			return errors.New("not implemented")
		}),
	}
}
