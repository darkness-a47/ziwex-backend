package types

import "github.com/labstack/echo/v4"

type Response struct {
	response   interface{}
	statusCode int
	err        error
}

func (r *Response) Write(status int, response interface{}) {
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
	c.JSON(r.statusCode, r.response)
	return nil
}
