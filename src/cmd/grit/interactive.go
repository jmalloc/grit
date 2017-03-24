package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli"
)

// choose asks the user to select an entry from opts interactively.
func choose(c *cli.Context, opt []string) (int, bool) {
	size := len(opt)

	if size == 0 {
		return 0, false
	} else if size == 1 {
		return 0, true
	}

	width := len(strconv.Itoa(size))
	f := fmt.Sprintf("  %%%dd) %%s", width)

	for i, o := range opt {
		write(c, f, i+1, o)
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Fprint(c.App.Writer, "> ")

		scanner.Scan()
		input := scanner.Text()

		switch strings.ToLower(input) {
		case "q", "quit":
			return 0, false
		default:
			i64, _ := strconv.ParseUint(input, 10, 64)
			idx := int(i64)

			if idx >= 1 && idx <= size {
				return idx - 1, true
			}
		}
	}
}
