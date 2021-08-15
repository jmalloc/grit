package shell

import (
	"io"
)

// Executor is a function used to execute commands within the context of Grit's
// parent shell.
type Executor func(command string, args ...string) error

// NewExecutor returns a new executor that writes commands to be executed to w.
func NewExecutor(w io.Writer) Executor {
	return func(command string, args ...string) error {
		if _, err := io.WriteString(w, Escape(command)); err != nil {
			return err
		}

		for _, arg := range args {
			if _, err := io.WriteString(w, ` `+Escape(arg)); err != nil {
				return err
			}
		}

		_, err := io.WriteString(w, "\n")
		return err
	}
}
