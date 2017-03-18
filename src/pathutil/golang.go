package pathutil

import (
	"os"
	"path"
)

// GoPath returns the current user's $GOPATH.
func GoPath() (string, error) {
	dir := os.Getenv("GOPATH")
	if dir != "" {
		return dir, nil
	}

	home, err := HomeDir()
	if err != nil {
		return "", err
	}

	return path.Join(home, "go"), nil
}
