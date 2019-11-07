package pathutil

import (
	"path/filepath"
	"strings"
)

// RelChild returns the relative path from base to p, if p is a child of base.
func RelChild(base, p string) (string, bool) {
	rel, err := filepath.Rel(base, p)
	if err != nil {
		return "", false
	}

	if strings.HasPrefix(rel, "..") {
		return "", false
	}

	return rel, true
}
