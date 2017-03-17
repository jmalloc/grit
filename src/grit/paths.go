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

	return path.Join(HomeDir(), "go"), true
}

// HomeDir returns the current user's home directory.
func HomeDir() string {
	dir := os.Getenv("HOME")
	if dir != "" {
		return dir
	}

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	if usr.HomeDir == "" {
		panic("could not determine home directory")
	}

	return usr.HomeDir
}
