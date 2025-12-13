package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type Server struct {
	host   string
	port   int32
	routes *echo.Echo
}

func NewServer(host string, port int32, routes *echo.Echo) Server {
	return Server{host, port, routes}
}

func (s Server) Start() {
	srv := http.Server{
		Addr:         fmt.Sprintf("%v:%v", s.host, s.port),
		Handler:      s.routes,
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
	}

	slog.Info("starting the server", "host", s.host, "port", s.port)
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
