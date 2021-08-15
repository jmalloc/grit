package ast

import "strconv"

// unquote returns the string value of a string token.
func unquote(s string) string {
	s, err := strconv.Unquote(s)
	if err != nil {
		panic(err)
	}

	return s
}
