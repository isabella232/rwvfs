// Package rwvfs augments vfs to support write operations.
package rwvfs

import (
	"io"
	"os"
	"syscall"

	"code.google.com/p/go.tools/godoc/vfs"
)

type FileSystem interface {
	vfs.FileSystem

	// Create creates the named file, truncating it if it already exists.
	Create(path string) (io.WriteCloser, error)

	// Mkdir creates a new directory. If name is already a directory, Mkdir
	// returns an error (that can be detected using os.IsExist).
	Mkdir(name string) error

	// Remove removes the named file or directory.
	Remove(name string) error
}

// MkdirAll creates a directory named path, along with any necessary parents. If
// path is already a directory, MkdirAll does nothing and returns nil.
func MkdirAll(fs FileSystem, path string) error {
	// adapted from os/MkdirAll

	dir, err := fs.Stat(path)
	if err == nil {
		if dir.IsDir() {
			return nil
		}
		return &os.PathError{"mkdir", path, syscall.ENOTDIR}
	}

	i := len(path)
	for i > 0 && os.IsPathSeparator(path[i-1]) {
		i--
	}

	j := i
	for j > 0 && !os.IsPathSeparator(path[j-1]) {
		j--
	}

	if j > 1 {
		err = MkdirAll(fs, path[0:j-1])
		if err != nil {
			return err
		}
	}

	err = fs.Mkdir(path)
	if err != nil {
		dir, err1 := fs.Lstat(path)
		if err1 == nil && dir.IsDir() {
			return nil
		}
		return err
	}
	return nil
}
