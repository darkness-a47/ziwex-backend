package controllers

import (
	"ziwex/dtos"
	"ziwex/services"

	"github.com/labstack/echo/v4"
)

func CreateProduct(c echo.Context) error {
	d := dtos.CreateProduct{}
	if err := c.Bind(&d); err != nil {
		return err
	}
	if err := c.Validate(&d); err != nil {
		return err
	}
	r := services.CreateProduct(d)
	return r.SendResponse(c)
}
