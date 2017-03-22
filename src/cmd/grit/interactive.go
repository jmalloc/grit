package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// choose asks the user to select an entry from opts interactively.
func choose(w io.Writer, opt []string) (int, bool) {
	if len(opt) == 0 {
		return 0, false
	} else if len(opt) == 1 {
		return 0, true
	}

	width := len(strconv.Itoa(len(opt)))
	f := fmt.Sprintf("  %%%dd) %%s\n", width)

	for i, o := range opt {
		fmt.Fprintf(w, f, i+1, o)
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Fprint(w, "> ")

		scanner.Scan()
		input := scanner.Text()

		switch strings.ToLower(input) {
		case "q", "quit":
			return 0, false
		default:
			i64, _ := strconv.ParseUint(input, 10, 64)
			idx := int(i64)

			if idx >= 1 && idx <= len(opt) {
				return idx - 1, true
			}
		}
	}
}
