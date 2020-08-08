package util

import (
	"os"
)

// IsDir path
func IsDir(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}
