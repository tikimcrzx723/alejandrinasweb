package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/tikimcrzx723/alejandrinasweb/views"
)

func Home(c echo.Context) error {
	// return views.
	return views.HomePage().Render(c.Request().Context(), c.Response())
}

func Product(c echo.Context) error {
	return views.ProductPage().Render(c.Request().Context(), c.Response())
}
