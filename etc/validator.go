package etc

import (
	"net/http"
	"strings"
	"ziwex/types/jsonResponse"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		//TODO: remove in prod
		e := strings.Split(err.Error(), "\n")
		return echo.NewHTTPError(http.StatusBadRequest, jsonResponse.Json{
			"message": e,
		})
	}
	return nil
}
