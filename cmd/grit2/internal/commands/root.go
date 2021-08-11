package commands

import (
	"embed"
	"strings"

	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "grit2",
	Short: "keep track of your local git clones",
}

//go:embed *.help
var helpFS embed.FS

// Root returns the root command.
func Root(version string) *cobra.Command {
	root.Version = version

	for _, cmd := range root.Commands() {
		if data, err := helpFS.ReadFile(cmd.Name() + ".help"); err == nil {
			const prefix = "  "
			help := prefix + strings.ReplaceAll(
				strings.TrimSpace(string(data)),
				"\n",
				"\n"+prefix,
			)
			cmd.Long = "Description:\n" + help
		}
	}

	return root
}
