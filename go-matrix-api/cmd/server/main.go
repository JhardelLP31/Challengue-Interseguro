// Punto de entrada de la API HTTP en Go + Fiber para operaciones matriciales.
package main

import (
	"log"

	"github.com/JhardelLP31/Challengue-Interseguro/go-matrix-api/internal/app"
	"github.com/JhardelLP31/Challengue-Interseguro/go-matrix-api/internal/client"
	"github.com/JhardelLP31/Challengue-Interseguro/go-matrix-api/internal/config"
)

func main() {
	cfg := config.Load()
	statsClient := client.NewHTTPStatsClient(cfg.StatsAPIURL, cfg.HTTPTimeout)
	fiberApp := app.New(cfg, statsClient)

	log.Printf("go-matrix-api escuchando en :%s (stats API: %s)", cfg.Port, cfg.StatsAPIURL)
	if err := fiberApp.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("error iniciando servidor Fiber: %v", err)
	}
}
