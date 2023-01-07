package routes

import (
	"ziwex/controllers"

	"github.com/labstack/echo/v4"
)

func AuthRoutesInit(g *echo.Group) {
	g.POST("/adminRegister", controllers.AuthAdminRegister)
	g.POST("/adminLogin", controllers.AuthAdminLogin)
}
