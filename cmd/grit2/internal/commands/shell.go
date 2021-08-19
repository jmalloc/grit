package commands

import (
	"errors"
	"io"
	"os"

	"github.com/jmalloc/grit/cmd/grit2/internal/di"
	"github.com/jmalloc/grit/internal/shell"
	"github.com/spf13/cobra"
)

// newSHellIntegrationCommand returns the "shell-integration" command.
func newShellIntegrationCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "shell-integration",
		Short: "setup shell integration",
		RunE: func(
			cmd *cobra.Command,
			args []string,
		) error {
			return errors.New("not implemented")
		},
	}
}

// setupShellExecutor adds the --shell-executor-output flag as a persistent flag
// on the root command so that it is available to all commands.
//
// It is marked as hidden as it should only be passed by the auto generated grit
// shell function, and never by the user directly.
func setupShellExecutor(root *cobra.Command) {
	f := root.PersistentFlags()
	f.String("shell-executor-output", "", "output file for shell commands to execute")
	f.MarkHidden("shell-executor-output") //nolint:errcheck
}

// provideShellExecutor adds a shell.Executor to the DI configuration.
func provideShellExecutor(cmd *cobra.Command) {
	di.Provide(func(d *di.Deferrer) (shell.Executor, error) {
		filename, err := cmd.Flags().GetString("shell-executor-output")
		if err != nil {
			return nil, err
		}

		if filename == "" {
			cmd.PrintErrf("Shell integration has not been configured. For more information run:\n\n")
			cmd.PrintErrf("    %s help shell-integration\n\n", executableName())
			return shell.NewExecutor(io.Discard), nil
		}

		fp, err := os.Create(filename)
		if err != nil {
			return nil, err
		}

		d.Defer(func() error {
			defer fp.Close()

			if err := fp.Sync(); err != nil {
				return err
			}

			return fp.Close()
		})

		return shell.NewExecutor(fp), nil
	})
}
