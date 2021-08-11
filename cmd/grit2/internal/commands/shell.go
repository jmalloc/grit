package commands

import (
	"errors"
	"io"
	"os"

	"github.com/jmalloc/grit/cmd/grit2/internal/di"
	"github.com/jmalloc/grit/shell"
	"github.com/spf13/cobra"
)

// shellExecutorOutputFlag is the name of the flag used to configure where shell
// commands are written.
const shellExecutorOutputFlag = "shell-executor-output"

// init adds the "shell-integration" command to the root command.
func init() {
	cmd := &cobra.Command{
		Use:   "shell-integration",
		Short: "setup shell integration",
		RunE: func(
			cmd *cobra.Command,
			args []string,
		) error {
			return errors.New("not implemented")
		},
	}

	root.AddCommand(cmd)

	// Add the shellExecutorOutputFlag as a persistent flag on the root command
	// so that it is available to all commands. It is marked as hidden as it
	// should only be passed by the auto generated grit shell function, and
	// never by the user directly.
	f := root.PersistentFlags()
	f.String(shellExecutorOutputFlag, "", "output file for shell commands to execute")
	f.MarkHidden(shellExecutorOutputFlag) //nolint:errcheck
}

// provideShellExecutor adds a shell.Executor to the DI configuration.
func provideShellExecutor(cmd *cobra.Command) error {
	return di.Provide(func(d *di.Deferrer) (shell.Executor, error) {
		filename, err := cmd.Flags().GetString(shellExecutorOutputFlag)
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
