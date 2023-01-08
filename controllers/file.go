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
	//TODO: check files mimetype
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

func GetFiles(c echo.Context) error {
	d := dtos.GetFiles{}
	if err := c.Bind(&d); err != nil {
		return err
	}
	if err := c.Validate(&d); err != nil {
		return err
	}
	r := services.GetFiles(d)
	return r.SendResponse(c)
}

func UpdateFile(c echo.Context) error {
	d := dtos.UpdateFile{}
	if err := c.Bind(&d); err != nil {
		return err
	}
	if err := c.Validate(&d); err != nil {
		return err
	}
	file, err := c.FormFile("file")
	if err != nil {
		d.File = nil
	} else {
		d.File = file
	}
	if d.Filename == nil && d.File == nil {
		r := jsonResponse.Response{}
		r.Write(http.StatusBadRequest, jsonResponse.Json{
			"message": "you must at least update one field",
		})
		return r.SendResponse(c)
	}
	r := services.UpdateFile(d)
	return r.SendResponse(c)
}

func DeleteFile(c echo.Context) error {
	d := dtos.DeleteFile{}
	if err := c.Bind(&d); err != nil {
		return err
	}
	if err := c.Validate(&d); err != nil {
		return err
	}
	r := services.DeleteFile(d)
	return r.SendResponse(c)
}
