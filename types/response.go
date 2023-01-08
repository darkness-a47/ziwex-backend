package types

import "github.com/labstack/echo/v4"

type Response interface {
	SendResponse(c echo.Context) error
}
