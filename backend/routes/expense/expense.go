package expense

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Aman-s12345/expense/backend/providers"
)

func RegisterRoutes(router fiber.Router, path string, prv *providers.Provider) {
	group := router.Group(path)

	group.Post("/", CreateExpenseHandler(prv))
	group.Get("/", ListExpensesHandler(prv))
}