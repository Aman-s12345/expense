package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func RateLimiter(max int, window time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: window,
		KeyGenerator: func(c *fiber.Ctx) string {
			return GetClientIP(c)
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"message": "Too many requests, please try again later",
			})
		},
		SkipFailedRequests: false,
	})
}

func StrictRateLimiter() fiber.Handler {
	return RateLimiter(10, 1*time.Minute)
}

func GeneralRateLimiter() fiber.Handler {
	return RateLimiter(60, 1*time.Minute)
}