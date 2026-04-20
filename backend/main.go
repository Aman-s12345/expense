package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/Aman-s12345/expense/backend/config"
	"github.com/Aman-s12345/expense/backend/providers"
	"github.com/Aman-s12345/expense/backend/routes"
)

func main() {
	cfg := config.NewConfig()

	log.SetLevel(cfg.Logger.Level)

	app := fiber.New(fiber.Config{
		AppName:        cfg.Server.Name,
		ErrorHandler:   customErrorHandler,
		ReadBufferSize: 8192,
	})

	setupMiddleware(app, cfg)

	provider, err := providers.NewProvider(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize providers: %v", err)
	}
	defer provider.Cleanup()

	app.Use(providers.Handle(provider))

	routes.RegisterRoutes(app, provider)

	// Log registered routes
	for _, route := range app.GetRoutes() {
		if route.Method == "OPTIONS" || route.Method == "HEAD" || route.Method == "TRACE" || route.Method == "CONNECT" {
			continue
		}
		if route.Path == "/" {
			continue
		}
		log.Infof("%s %s", route.Method, route.Path)
	}

	address := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Infow("Starting server", "address", address)

	if err := app.Listen(address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupMiddleware(app *fiber.App, cfg *config.Config) {
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.Server.FrontendURL,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": message,
		"error":   err.Error(),
	})
}