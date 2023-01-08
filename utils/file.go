package utils

import (
	"io"
	"net/http"
)

func GetFileContentType(out io.ReadSeeker) (string, error) {

	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)

	out.Seek(0, 0)
	return contentType, nil
}
