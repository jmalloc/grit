package main

import (
	"context"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmalloc/grit/cmd/grit2/internal/commands"
	"github.com/jmalloc/grit/cmd/grit2/internal/di"
	"go.uber.org/multierr"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	if err := run(); err != nil {
		os.Exit(1)
	}
}

// version string, automatically set during build process.
var version = "0.0.0"

func run() (err error) {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	defer func() {
		err = multierr.Append(
			err,
			di.Close(),
		)
	}()

	return commands.Root(version).ExecuteContext(ctx)
}
