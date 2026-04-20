package auth

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Aman-s12345/expense/backend/middlewares"
	"github.com/Aman-s12345/expense/backend/providers"
)

func RegisterPublicRoutes(router fiber.Router, path string, prv *providers.Provider) {
	group := router.Group(path, middlewares.StrictRateLimiter())

	group.Post("/register", RegisterHandler(prv))
	group.Post("/login", LoginHandler(prv))
	group.Post("/guest", GuestLoginHandler(prv))
}

func RegisterProtectedRoutes(router fiber.Router, path string, prv *providers.Provider) {
	group := router.Group(path)

	group.Post("/logout", LogoutHandler(prv))
	group.Get("/me", MeHandler(prv))
}