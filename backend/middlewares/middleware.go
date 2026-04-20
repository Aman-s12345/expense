package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"

	"github.com/Aman-s12345/expense/backend/providers"
)

func RequireAuth(prv *providers.Provider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			log.Warn("Missing authorization header")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Authorization header is required",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			log.Warn("Invalid token format")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Invalid authorization format. Use: Bearer <token>",
			})
		}

		user, err := prv.S.Auth.ValidateSession(c.Context(), tokenString)
		if err != nil {
			log.Warnw("Session validation failed", "error", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Invalid or expired session",
			})
		}

		c.Locals("user_id", user.ID.String())
		c.Locals("user_uuid", user.ID)
		c.Locals("is_guest", user.IsGuest)

		log.Debugw("User authenticated", "user_id", user.ID, "is_guest", user.IsGuest)
		return c.Next()
	}
}

func GetUserID(c *fiber.Ctx) (uuid.UUID, error) {
	userUUID, ok := c.Locals("user_uuid").(uuid.UUID)
	if !ok {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, "user not authenticated")
	}
	return userUUID, nil
}

func GetClientIP(c *fiber.Ctx) string {
	if realIP := c.Get("X-Real-IP"); realIP != "" {
		return strings.TrimSpace(realIP)
	}

	if forwardedFor := c.Get("X-Forwarded-For"); forwardedFor != "" {
		ips := strings.Split(forwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	return c.IP()
}