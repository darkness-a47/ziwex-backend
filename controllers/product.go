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

func GetProductsSummery(c echo.Context) error {
	d := dtos.GetProductsSummery{}
	if err := c.Bind(&d); err != nil {
		return err
	}
	if err := c.Validate(&d); err != nil {
		return err
	}
	url := c.Request().URL.RequestURI()
	d.RequestPath = url
	r := services.GetProductsSummery(d)
	return r.SendResponse(c)
}

func GetProductData(c echo.Context) error {
	d := dtos.GetProductData{}
	if err := c.Bind(&d); err != nil {
		return err
	}
	if err := c.Validate(&d); err != nil {
		return err
	}
	url := c.Request().URL.RequestURI()
	d.RequestPath = url
	r := services.GetProductData(d)
	return r.SendResponse(c)
}
