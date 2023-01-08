package fileResponse

import (
	"fmt"
	"io"

	"github.com/labstack/echo/v4"
)

type File struct {
	File           io.ReadCloser
	ContentType    string
	CancelFunction func()
}

type Response struct {
	response   *File
	statusCode int
	err        error
}

func (r *Response) Write(status int, response *File) {
	r.statusCode = status
	r.response = response
	r.err = nil
}

func (r *Response) Error(err error) {
	r.response = nil
	r.statusCode = 500
	r.err = err
}

func (r *Response) SendResponse(c echo.Context) error {
	if r.err != nil {
		return r.err
	}
	if r.response != nil {
		c.Stream(r.statusCode, r.response.ContentType, r.response.File)
		defer r.response.File.Close()
		defer r.response.CancelFunction()
		return nil
	}
	return fmt.Errorf("no response written")
}
