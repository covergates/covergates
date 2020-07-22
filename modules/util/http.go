package util

import (
	"bytes"
	"mime/multipart"
	"net/http"
)

// FormFile defines file name and data for filed
type FormFile struct {
	Name string
	Data []byte
}

// FormData defines form field for http POST
type FormData map[string]interface{}

// PostForm request
func PostForm(url string, form FormData) (*http.Response, error) {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	for key, value := range form {
		switch v := value.(type) {
		case string:
			if err := w.WriteField(key, v); err != nil {
				return nil, err
			}
		case FormFile:
			file, err := w.CreateFormFile(key, v.Name)
			if err != nil {
				return nil, err
			}
			if _, err := file.Write(v.Data); err != nil {
				return nil, err
			}
		}
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	c := client()
	return c.Do(req)
}

func client() *http.Client {
	return new(http.Client)
}
