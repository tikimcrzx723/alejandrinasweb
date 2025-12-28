package contexts

import (
	"github.com/labstack/echo/v4"
)

type AppKey struct{}

func (AppKey) String() string {
	return "appCtx"
}

type App struct {
	echo.Context
	UserID          int
	IsAuthenticated bool
	Token           string
	Role            string
}
