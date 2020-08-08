package common

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/covergates/covergates/modules/util"
)

// OpenFileReader from path
func OpenFileReader(path string) (io.Reader, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return nil, fmt.Errorf("%s is not file", path)
	}
	file, err := os.Open(path)
	defer file.Close()
	buf := &bytes.Buffer{}
	io.Copy(buf, file)
	return buf, nil
}

// FindReport in path for target names
func FindReport(path string, target ...string) (string, error) {
	if !util.IsDir(path) {
		return path, nil
	}
	report := ""
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && searchStrings(target, filepath.Base(path)) {
			report = path
			return io.EOF
		}
		return nil
	})
	if err != nil && err != io.EOF {
		return "", err
	}
	if report == "" {
		return report, fmt.Errorf("report not found")
	}
	return report, nil
}

func searchStrings(slice []string, x string) bool {
	for _, s := range slice {
		if s == x {
			return true
		}
	}
	return false
}
