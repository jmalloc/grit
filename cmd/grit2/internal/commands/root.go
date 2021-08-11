package commands

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// root is the top-level "grit" command.
var root = &cobra.Command{
	Use:   executableName(),
	Short: "keep track of your local git clones",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return provideShellExecutor(cmd)
	},
}

// Root returns the root command.
//
// v is the version to display. It is passed from the main package where it is
// made available as part of the build process.
func Root(v string) *cobra.Command {
	setupHelp()
	root.Version = v

	return root
}

// executableName returns the name of the grit executable.
func executableName() string {
	return filepath.Base(os.Args[0])
}
