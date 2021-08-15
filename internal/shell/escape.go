package shell

import (
	"strings"
)

// Escape quotes and escapes a string for use as a shell argument.
func Escape(s string) string {
	return `'` + strings.ReplaceAll(s, `'`, `'"'"'`) + `'`
}
