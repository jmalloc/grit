package autocomplete

import (
	"fmt"
	"os"
	"strings"

	"github.com/jmalloc/grit/src/grit"
	"github.com/jmalloc/grit/src/grit/index"
	"github.com/urfave/cli"
)

// Function provides auto-completion for a specific argument.
type Function func(
	cfg grit.Config,
	idx *index.Index,
	c *cli.Context,
	arg string,
) []string

// New returns a autocompletion functions that resolves arguments in order.
func New(funcs ...Function) cli.BashCompleteFunc {
	return func(c *cli.Context) {
		i := c.NArg() - 1
		if strings.HasSuffix(os.Getenv("GRIT_COMP_WORDS"), " ") {
			i++
		}

		if i < 0 || i >= len(funcs) {
			return
		}

		arg := ""
		if i < len(c.Args()) {
			arg = c.Args()[i]
		}

		cfg, err := grit.LoadConfig(c.GlobalString("config"))
		if err != nil {
			panic(err)
		}

		idx, err := index.Open(cfg)
		if err != nil {
			panic(err)
		}
		defer idx.Close()

		fn := funcs[i]
		for _, str := range fn(cfg, idx, c, arg) {
			fmt.Fprintln(c.App.Writer, str)
		}
	}
}
