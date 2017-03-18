package pathutil

import (
	"os"
	"path"
	"strings"
)

// Resolve resolves p to an absolute path.
func Resolve(p string) (string, error) {
	if path.IsAbs(p) {
		return p, nil
	}

	base, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return ResolveFrom(base, p)
}

// ResolveFrom resolves p to an absolute path, relative to base.
func ResolveFrom(base, p string) (string, error) {
	if p == "" {
		return base, nil
	} else if path.IsAbs(p) {
		return p, nil
	}

	pos := strings.IndexByte(p, os.PathSeparator)
	var head, tail string
	if pos == -1 {
		head = p
	} else {
		head = p[:pos]
		tail = p[pos+1:]
	}

	if head == "~" {
		home, err := HomeDir()
		if err != nil {
			return "", err
		}
		return path.Join(home, tail), nil
	}

	if head[0] == '~' {
		home, err := HomeDirOf(head[1:])
		if err != nil {
			return "", err
		}
		return path.Join(home, tail), nil
	}

	return path.Join(base, p), nil
}
