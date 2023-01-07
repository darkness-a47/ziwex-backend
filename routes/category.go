package routes

import (
	"ziwex/controllers"
	"ziwex/middlewares"

	"github.com/labstack/echo/v4"
)

func CategoryRoutesInit(g *echo.Group) {
	g.POST("/create", controllers.CreateCategory, middlewares.AuthorizeAccess("admin"))
	g.GET("/categories", controllers.GetCategories, middlewares.AuthorizeAccess("admin", "user"))
}
