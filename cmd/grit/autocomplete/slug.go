package autocomplete

import (
	"github.com/jmalloc/grit"
	"github.com/jmalloc/grit/src/grit/index"
	"github.com/urfave/cli"
)

// Slug provides auto-completion for indexed slugs.
func Slug(cfg grit.Config, idx *index.Index, c *cli.Context, arg string) []string {
	return idx.ListSlugs(arg)
}
