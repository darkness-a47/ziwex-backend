package emptyResponse

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type Response struct {
	statusCode int
	err        error
}

func (r *Response) Write(status int) {
	r.statusCode = status
	r.err = nil
}

func (r *Response) Error(err error) {
	r.statusCode = 500
	r.err = err
}

func (r *Response) SendResponse(c echo.Context) error {
	if r.err != nil {
		return r.err
	}
	if r.statusCode != 0 {
		c.NoContent(r.statusCode)
		return nil
	}
	return fmt.Errorf("no response written")
}
