package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func promptBetween(min, max int) int {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Fprintf(os.Stderr, "> ")

		scanner.Scan()
		input := scanner.Text()

		i64, _ := strconv.ParseUint(input, 10, 64)
		i := int(i64)

		if i >= min && i <= max {
			return i
		}
	}
}
