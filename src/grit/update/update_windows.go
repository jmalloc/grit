// +build windows

package update

import "runtime"

const (
	archiveName       = "grit-windows-" + runtime.GOARCH + ".zip"
	archiveBinaryName = "grit.exe"
)

// Unpack an archive and return the location of the inner binary.
func Unpack(src, dst string) error {
	panic("not implemented")
}
