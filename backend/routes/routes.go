package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Aman-s12345/expense/backend/middlewares"
	"github.com/Aman-s12345/expense/backend/providers"
	authRoutes "github.com/Aman-s12345/expense/backend/routes/auth"
	expenseRoutes "github.com/Aman-s12345/expense/backend/routes/expense"
)

func RegisterRoutes(app *fiber.App, prv *providers.Provider) {
	api := app.Group("/api")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "OK",
		})
	})

	public := api.Group("", middlewares.GeneralRateLimiter())
	authRoutes.RegisterPublicRoutes(public, "/auth", prv)

	protected := api.Group("", middlewares.GeneralRateLimiter(), middlewares.RequireAuth(prv))
	authRoutes.RegisterProtectedRoutes(protected, "/auth", prv)
	expenseRoutes.RegisterRoutes(protected, "/expenses", prv)
}