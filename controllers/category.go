package controllers

import (
	"ziwex/dtos"
	"ziwex/services"

	"github.com/labstack/echo/v4"
)

func CreateCategory(c echo.Context) error {
	d := dtos.CreateCategory{}
	if err := c.Bind(&d); err != nil {
		return err
	}
	if err := c.Validate(&d); err != nil {
		return err
	}
	r := services.CreateCategory(d)
	return r.SendResponse(c)
}

func GetCategories(c echo.Context) error {
	d := dtos.GetCategories{}
	if err := c.Bind(&d); err != nil {
		return err
	}
	if err := c.Validate(&d); err != nil {
		return err
	}
	r := services.GetCategories(d)
	return r.SendResponse(c)
}
