package main

import (
	"github.com/tikimcrzx723/alejandrinasweb/internal/env"
	"github.com/tikimcrzx723/alejandrinasweb/routes"
	"github.com/tikimcrzx723/alejandrinasweb/server"
)

func main() {
	routes := routes.NewRoutes()
	host := env.GetString("SERVER_HOST", "0.0.0.0")
	port := env.GetInt("SERVER_PORT", 8080)

	srv := server.NewServer(host, int32(port), routes.Load())

	srv.Start()
}
