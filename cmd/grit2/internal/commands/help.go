package commands

import (
	"embed"
	"strings"
)

// helpFS is an embedded filesystem containing the help content for each
// command.
//
//go:embed *.txt
var helpFS embed.FS

// setupHelp configures each command within the root command to use its
// associated .txt file as its "long message", which is displayed by the "help"
// sub-command.
func setupHelp() {
	for _, cmd := range root.Commands() {
		if data, err := helpFS.ReadFile(cmd.Name() + ".txt"); err == nil {
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
}
