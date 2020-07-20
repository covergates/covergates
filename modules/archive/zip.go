package archive

import (
	"archive/zip"
	"bytes"
	"io"
)

// NewZipReader from io.Reader
func NewZipReader(r io.Reader) (*zip.Reader, error) {
	buffer := bytes.NewBuffer([]byte{})
	size, err := io.Copy(buffer, r)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(buffer.Bytes())
	return zip.NewReader(reader, size)
}
