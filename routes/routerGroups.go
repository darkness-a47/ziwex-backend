package routes

import "github.com/labstack/echo/v4"

func RouterGroupsInit(e *echo.Echo) {
	AuthRoutesInit(e.Group("/auth"))
	CategoryRoutesInit(e.Group("/category"))
	FileRoutesInit(e.Group("/file"))
	ProductRoutesInit(e.Group("/product"))
}
