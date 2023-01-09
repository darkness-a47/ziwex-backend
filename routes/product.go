package routes

import (
	"ziwex/controllers"

	"github.com/labstack/echo/v4"
)

func ProductRoutesInit(g *echo.Group) {
	g.POST("/create", controllers.CreateProduct)
}
