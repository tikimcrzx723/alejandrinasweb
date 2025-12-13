package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/tikimcrzx723/alejandrinasweb/controllers"
	"github.com/tikimcrzx723/alejandrinasweb/static"
)

type Routes struct {
	e *echo.Echo
}

func NewRoutes() Routes {
	e := echo.New()
	echo.MustSubFS(static.Files, "static")
	e.StaticFS("/static", static.Files)
	return Routes{e}
}

func (r Routes) Load() *echo.Echo {
	// setup routes for diferents pages
	r.e.GET("", func(c echo.Context) error {
		return controllers.Home(c)
	})
	r.e.GET("/product", func(c echo.Context) error {
		return controllers.Product(c)
	})
	return r.e
}
