package commands

import (
	"embed"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Short: "keep track of your local git clones",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return provideShellExecutor(cmd)
	},
}

//go:embed *.help
var helpFS embed.FS

// Root returns the root command.
func Root(version string) *cobra.Command {
	root.Use = executableName()
	root.Version = version

	for _, cmd := range root.Commands() {
		if data, err := helpFS.ReadFile(cmd.Name() + ".help"); err == nil {
			const prefix = "  "

			help := strings.ReplaceAll(
				strings.TrimSpace(string(data)),
				"{{ executable }}",
				executableName(),
			)

			cmd.Long = "Description:\n" +
				prefix +
				strings.ReplaceAll(
					help,
					"\n",
					"\n"+prefix,
				)
		}
	}

	return root
}

// executableName returns the name of the grit executable.
func executableName() string {
	return filepath.Base(os.Args[0])
}
