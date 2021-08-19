package commands

import (
	"os"
	"path/filepath"

	"github.com/jmalloc/grit/cmd/grit2/internal/commands/source"
	"github.com/spf13/cobra"
)

// NewRoot returns the root command.
//
// v is the version to display. It is passed from the main package where it is
// made available as part of the build process.
func NewRoot(v string) *cobra.Command {
	var root *cobra.Command
	root = &cobra.Command{
		Version: v,
		Use:     executableName(),
		Short:   "keep track of your local git clones",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			provideConfig(cmd)
			provideShellExecutor(cmd)
		},
	}

	setupConfig(root)
	setupShellExecutor(root)

	root.AddCommand(
		source.NewRoot(),
		newShellIntegrationCommand(),
		newCloneCommand(),
	)

	return root
}

// executableName returns the name of the grit executable.
func executableName() string {
	return filepath.Base(os.Args[0])
}
