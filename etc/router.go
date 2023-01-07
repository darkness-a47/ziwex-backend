package etc

import (
	"ziwex/routes"

	"github.com/labstack/echo/v4"
)

func RouterInit(e *echo.Echo) {
	routes.RouterGroupsInit(e)
}
