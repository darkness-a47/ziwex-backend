package routes

import (
	"ziwex/controllers"

	"github.com/labstack/echo/v4"
)

func ProductRoutesInit(g *echo.Group) {
	g.POST("/create", controllers.CreateProduct)
	g.GET("/productsSummery", controllers.GetProductsSummery)
	g.GET("/data/:product_url", controllers.GetProductData)
}
