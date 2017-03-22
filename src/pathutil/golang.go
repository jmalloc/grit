package pathutil

import (
	"os"
	"path"
)

// GoPath returns the current user's $GOPATH directory.
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

// GoSrc returns the current user's $GOPATH/src directory.
func GoSrc() (string, error) {
	p, ok := GoPath()
	return path.Join(p, "src"), ok
}
