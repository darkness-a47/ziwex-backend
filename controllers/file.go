package controllers

import (
	"net/http"
	"ziwex/dtos"
	"ziwex/services"
	"ziwex/types/jsonResponse"

	"github.com/labstack/echo/v4"
)

func UploadFile(c echo.Context) error {
	d := dtos.UploadFile{}
	if err := c.Bind(&d); err != nil {
		return err
	}
	if err := c.Validate(&d); err != nil {
		return err
	}
	file, err := c.FormFile("file")
	if err != nil {
		r := jsonResponse.Response{}
		r.Write(http.StatusBadRequest, jsonResponse.Json{
			"message": "file not found",
		})
		return r.SendResponse(c)
	}

	d.File = file
	res := services.UploadFile(d)
	return res.SendResponse(c)
}

func ServeFile(c echo.Context) error {
	d := dtos.ServeFile{}
	if err := c.Bind(&d); err != nil {
		return err
	}
	if err := c.Validate(&d); err != nil {
		return err
	}
	res := services.ServeFile(d)
	return res.SendResponse(c)
}
