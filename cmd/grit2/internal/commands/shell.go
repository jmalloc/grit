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

	// Add the
	pflags := root.PersistentFlags()
	pflags.String(shellExecutorOutputFlag, "", "")
	if err := pflags.MarkHidden(shellExecutorOutputFlag); err != nil {
		panic(err)
	}
}

// provideShellExecutor adds a shell.Executor to the DI configuration.
func provideShellExecutor(cmd *cobra.Command) error {
	return di.Provide(cmd, func(d *di.Deferrer) (shell.Executor, error) {
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
