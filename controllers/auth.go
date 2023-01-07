package controllers

import (
	"ziwex/dtos"
	"ziwex/services"

	"github.com/labstack/echo/v4"
)

func AuthAdminRegister(c echo.Context) error {
	user := dtos.AuthAdminRegister{}
	if err := c.Bind(&user); err != nil {
		return err
	}
	if err := c.Validate(&user); err != nil {
		return err
	}
	r := services.AuthAdminRegister(user)
	return r.SendResponse(c)
}

func AuthAdminLogin(c echo.Context) error {
	user := dtos.AuthAdminLogin{}
	if err := c.Bind(&user); err != nil {
		return err
	}
	if err := c.Validate(&user); err != nil {
		return err
	}
	r := services.AuthAdminLogin(user)
	return r.SendResponse(c)
}
