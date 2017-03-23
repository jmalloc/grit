// +build !windows

package update

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"runtime"
)

const (
	archiveName       = "grit-" + runtime.GOOS + "-" + runtime.GOARCH + ".tar.gz"
	archiveBinaryName = "grit"
)

// Unpack an archive and return the location of the inner binary.
func Unpack(src, dst string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	z, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer z.Close()

	r := tar.NewReader(z)
	if err != nil {
		return err
	}

	for {
		header, err := r.Next()
		if err == io.EOF {
			return errors.New("could not find binary in archive")
		} else if err != nil {
			return err
		} else if header.Name != archiveBinaryName {
			continue
		}

		info := header.FileInfo()

		if info.IsDir() {
			return errors.New("unexpected directory in archive")
		}

		w, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer w.Close()

		_, err = io.Copy(w, r)
		return err
	}
}
