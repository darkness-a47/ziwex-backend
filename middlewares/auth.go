package middlewares

import (
	"net/http"
	"strings"
	"ziwex/types/jsonResponse"
	"ziwex/utils"

	"github.com/labstack/echo/v4"
)

func AuthorizeAccess(accessLevels ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			bearer := strings.Split(c.Request().Header.Get("Auth"), " ")

			r := jsonResponse.Response{}
			r.Write(http.StatusUnauthorized, jsonResponse.Json{
				"message": "unauthorized",
			})

			if len(bearer) != 2 || bearer[0] != "Bearer" {
				return r.SendResponse(c)
			}

			claims, err := utils.JwtValidateToken(bearer[1])
			if err != nil {
				return r.SendResponse(c)
			}
			claimsValue := *claims

			user := claimsValue["userType"].(string)
			if !utils.ListContainsString(accessLevels, user) {
				return r.SendResponse(c)
			}

			c.Set("username", claimsValue["username"])
			c.Set("userType", claimsValue["userType"])

			return next(c)
		}
	}
}
