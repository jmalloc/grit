//go:build windows
// +build windows

package update

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"runtime"
)

const (
	archiveName       = "grit-windows-" + runtime.GOARCH + ".zip"
	archiveBinaryName = "grit.exe"
)

// Unpack an archive and return the location of the inner binary.
func Unpack(src, dst string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name != archiveBinaryName {
			continue
		}

		info := f.FileInfo()

		if info.IsDir() {
			return errors.New("unexpected directory in archive")
		}

		w, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer w.Close()

		r, err := f.Open()
		if err != nil {
			return err
		}
		defer r.Close()

		_, err = io.Copy(w, r)
		return err
	}

	return errors.New("could not find binary in archive")
}
