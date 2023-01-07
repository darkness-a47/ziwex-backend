package main

import (
	"ziwex/db"
	"ziwex/etc"
	"ziwex/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	//TODO: remove in prod
	e.Debug = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = &etc.CustomValidator{Validator: validator.New()}

	close := db.ConnetionInit()
	defer close()

	etc.RouterInit(e)
	utils.JwtInit()

	e.Logger.Fatal(e.Start(":4200"))
}
