package middleware

import (
	"net/http"

	"github.com/tikimcrzx723/alejandrinasweb/controllers"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get(controllers.AuthSessionName, c)
		if err != nil {
			return c.Redirect(http.StatusTemporaryRedirect, "/login")
		}

		isAuth, _ := sess.Values[controllers.AuthUserAuthenticated].(bool)
		if isAuth {
			return next(c)
		}

		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
}

func RequireNoAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get(controllers.AuthSessionName, c)
		if err != nil {
			return next(c)
		}

		isAuth, _ := sess.Values[controllers.AuthUserAuthenticated].(bool)
		if isAuth {
			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}

		return next(c)
	}
}

func RequireAdminRole(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get(controllers.AuthSessionName, c)
		if err != nil {
			return c.Redirect(http.StatusTemporaryRedirect, "/login")
		}

		role, _ := sess.Values["ROLE"].(string)
		if role == "admin" {
			return next(c)
		}

		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
}
