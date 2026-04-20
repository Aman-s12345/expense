package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	authService "github.com/Aman-s12345/expense/backend/services/auth"
	"github.com/Aman-s12345/expense/backend/middlewares"
	"github.com/Aman-s12345/expense/backend/providers"
)

func RegisterHandler(prv *providers.Provider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req authService.RegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Invalid request body",
			})
		}

		user, session, err := prv.S.Auth.Register(c.Context(), req)
		if err != nil {
			log.Warnw("Registration failed", "email", req.Email, "error", err)

			status := fiber.StatusBadRequest
			if strings.Contains(err.Error(), "already registered") {
				status = fiber.StatusConflict
			}

			return c.Status(status).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		log.Infow("User registered", "user_id", user.ID, "email", req.Email)
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": true,
			"message": "Registration successful",
			"data": fiber.Map{
				"user":  user,
				"token": session.Token,
			},
		})
	}
}

func LoginHandler(prv *providers.Provider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req authService.LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Invalid request body",
			})
		}

		user, session, err := prv.S.Auth.Login(c.Context(), req)
		if err != nil {
			log.Warnw("Login failed", "email", req.Email, "error", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Invalid email or password",
			})
		}

		log.Infow("User logged in", "user_id", user.ID)
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Login successful",
			"data": fiber.Map{
				"user":  user,
				"token": session.Token,
			},
		})
	}
}

func GuestLoginHandler(prv *providers.Provider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, session, err := prv.S.Auth.GuestLogin(c.Context())
		if err != nil {
			log.Errorw("Guest login failed", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Failed to create guest session",
			})
		}

		log.Infow("Guest user created", "user_id", user.ID)
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": true,
			"message": "Guest session created",
			"data": fiber.Map{
				"user":  user,
				"token": session.Token,
			},
		})
	}
}

func LogoutHandler(prv *providers.Provider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		token := strings.TrimPrefix(authHeader, "Bearer ")

		if err := prv.S.Auth.Logout(c.Context(), token); err != nil {
			log.Warnw("Logout failed", "error", err)
			// Still return success — logout should be idempotent from the client's perspective
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Logged out successfully",
		})
	}
}

func MeHandler(prv *providers.Provider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := middlewares.GetUserID(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Not authenticated",
			})
		}

		user, err := prv.S.Auth.GetUserByID(c.Context(), userID)
		if err != nil {
			log.Warnw("Failed to get user", "user_id", userID, "error", err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "User not found",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    user,
		})
	}
}