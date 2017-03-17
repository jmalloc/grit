package grit

import (
	"os"
	"os/user"
	"path"
)

// GoPath returns the current user's Go path.
func GoPath() (string, bool) {
	dir := os.Getenv("GOPATH")
	if dir != "" {
		return dir, true
	}

	home, ok := HomeDir()
	if ok {
		return path.Join(home, "go"), true
	}

	return "", false
}

// HomeDir returns the current user's home directory.
func HomeDir() (string, bool) {
	dir := os.Getenv("HOME")
	if dir != "" {
		return dir, true
	}

	usr, err := user.Current()
	if err != nil || usr.HomeDir == "" {
		return "", false
	}

	return usr.HomeDir, true
}
