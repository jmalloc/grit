package commands

import "github.com/spf13/cobra"

var Root = &cobra.Command{
	Use:   "grit2",
	Short: "keep track of your local git clones",
}
