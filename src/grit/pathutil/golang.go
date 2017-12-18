package pathutil

import (
	"os"
	"path"
	"path/filepath"
)

// GoPath returns the current user's $GOPATH directory.
func GoPath() (string, error) {
	p := os.Getenv("GOPATH")
	if p != "" {
		dirs := filepath.SplitList(p)
		return dirs[0], nil
	}

	home, err := HomeDir()
	if err != nil {
		return "", err
	}

	return path.Join(home, "go"), nil
}

// GoSrc returns the current user's $GOPATH/src directory.
func GoSrc() (string, error) {
	p, err := GoPath()
	return path.Join(p, "src"), err
}
