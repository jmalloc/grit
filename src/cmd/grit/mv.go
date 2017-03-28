package main

import (
	"errors"
	"os"
	"path"

	"github.com/jmalloc/grit/src/grit"
	"github.com/jmalloc/grit/src/grit/index"
	"github.com/urfave/cli"
)

func mv(cfg grit.Config, idx *index.Index, c *cli.Context) error {
	src, err := dirFromFirstArg(c)
	if err != nil {
		return err
	}

	base, err := cloneBaseDir(cfg, c)
	if err != nil {
		return err
	}

	endpoints, err := grit.EndpointsFromDir(src)
	if err != nil {
		return err
	}

	var dirs []string
	for _, ep := range endpoints {
		dirs = append(dirs, grit.EndpointToDir(base, ep))
	}

	dst, ok := chooseCloneDir(cfg, c, dirs)
	if !ok {
		return errSilentFailure
	}

	write(c, dst)

	if src == dst {
		return nil
	}

	if c.NArg() == 0 {
		exec(c, "cd", dst)
	}

	_, err = os.Stat(dst)
	if err == nil {
		return errors.New("destination directory already exists")
	} else if !os.IsNotExist(err) {
		return err
	}

	if err := os.MkdirAll(path.Dir(dst), 0755); err != nil {
		return err
	}

	if err := os.Rename(src, dst); err != nil {
		return err
	}

	if err := idx.Add(dst); err != nil {
		return err
	}

	return idx.Remove(src)
}
