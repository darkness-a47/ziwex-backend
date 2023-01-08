package routes

import (
	"ziwex/controllers"

	"github.com/labstack/echo/v4"
)

func FileRoutesInit(g *echo.Group) {
	g.POST("/uploadFile", controllers.UploadFile)
	g.GET("/serve/:file_id", controllers.ServeFile)
	g.GET("/files", controllers.GetFiles)
}
