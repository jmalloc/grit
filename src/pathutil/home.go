package pathutil

import (
	"errors"
	"fmt"
	"os"
	"os/user"
)

// HomeDir returns the current user's home directory.
func HomeDir() (string, error) {
	dir := os.Getenv("HOME")
	if dir != "" {
		return dir, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	if usr.HomeDir == "" {
		return "", errors.New("home directory is not configured")
	}

	return usr.HomeDir, nil
}

// HomeDirOf returns a user's home directory.
func HomeDirOf(u string) (string, error) {
	usr, err := user.Lookup(u)
	if err != nil {
		return "", err
	}

	if usr.HomeDir == "" {
		return "", fmt.Errorf("home directory for %s is not configured", u)
	}

	return usr.HomeDir, nil
}
