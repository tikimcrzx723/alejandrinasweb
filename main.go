package main

import (
	"encoding/gob"

	"github.com/google/uuid"
	"github.com/tikimcrzx723/alejandrinasweb/internal/env"
	"github.com/tikimcrzx723/alejandrinasweb/routes"
	"github.com/tikimcrzx723/alejandrinasweb/routes/contexts"
	"github.com/tikimcrzx723/alejandrinasweb/server"
)

func main() {
	gob.Register(uuid.UUID{})
	gob.Register(contexts.FlashMessage{})

	routes := routes.NewRoutes()
	host := env.GetString("SERVER_HOST", "0.0.0.0")
	port := env.GetInt("SERVER_PORT", 9090)

	srv := server.NewServer(host, int32(port), routes.Load())

	srv.Start()
}
