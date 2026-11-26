package app

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/JhardelLP31/Challengue-Interseguro/go-matrix-api/internal/client"
	"github.com/JhardelLP31/Challengue-Interseguro/go-matrix-api/internal/config"
	"github.com/JhardelLP31/Challengue-Interseguro/go-matrix-api/internal/handler"
	"github.com/JhardelLP31/Challengue-Interseguro/go-matrix-api/internal/middleware"
	"github.com/JhardelLP31/Challengue-Interseguro/go-matrix-api/internal/model"
	"github.com/JhardelLP31/Challengue-Interseguro/go-matrix-api/internal/service"
)

// New construye y configura la aplicación Fiber con todas sus rutas y middlewares.
func New(cfg config.Config, statsClient client.StatsClient) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "go-matrix-api",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	})

	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(logger.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"service": "go-matrix-api",
		})
	})

	matrixService := service.NewMatrixService()
	authHandler := handler.NewAuthHandler(cfg)
	matrixHandler := handler.NewMatrixHandler(matrixService, statsClient)

	v1 := app.Group("/api/v1")
	v1.Post("/auth/token", authHandler.IssueToken)

	matrices := v1.Group("/matrices", middleware.JWTProtected(cfg.JWTSecret))
	matrices.Post("/qr", matrixHandler.FactorizeQR)
	matrices.Post("/rotate", matrixHandler.RotateMatrix)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
			Error:   "NotFound",
			Message: "Ruta no encontrada.",
		})
	})

	return app
}
