// package fsutil includes some util functions for operating file system
package fsutil

import (
	"errors"
	"io"
	"os"
)

// IsDir check if a path is a directory
func IsDir(pth string) (bool, error) {
	fi, err := os.Stat(pth)
	if err != nil {
		return false, err
	}

	return fi.IsDir(), nil
}

// FileExists check if the given file exists
func FileExists(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err == nil {
		if stat.Mode()&os.ModeType == 0 {
			return true, nil
		}
		return false, errors.New(path + " exists but is not regular file")
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// CopyFile copy src to dst
func CopyFile(src, dst string) (bool, error) {
	from, err := os.Open(src)
	if err != nil {
		return false, err
	}
	defer from.Close()

	to, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return false, err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return false, err
	}

	return true, nil
}
