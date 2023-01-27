package middlewares

import (
	"fmt"
	"net/http"
	"ziwex/cache"
	"ziwex/types/jsonbResponse"

	"github.com/labstack/echo/v4"
)

func CheckCacheParam() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			url := c.Request().URL.RequestURI()
			data, err := cache.Get(url, "")
			if err != nil {
				return next(c)
			}
			fmt.Println("cache used")
			r := jsonbResponse.Response{}
			r.Write(http.StatusOK, data)
			return r.SendResponse(c)
		}
	}
}
