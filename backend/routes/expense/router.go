package expense

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	"github.com/Aman-s12345/expense/backend/db/models"
	"github.com/Aman-s12345/expense/backend/middlewares"
	"github.com/Aman-s12345/expense/backend/providers"
	expenseService "github.com/Aman-s12345/expense/backend/services/expense"
)

func CreateExpenseHandler(prv *providers.Provider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := middlewares.GetUserID(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Not authenticated",
			})
		}

		var req expenseService.CreateExpenseRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Invalid request body",
			})
		}

		expense, err := prv.S.Expense.Create(c.Context(), userID, req)
		if err != nil {
			log.Warnw("Failed to create expense", "user_id", userID, "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": true,
			"message": "Expense created",
			"data":    toExpenseResponse(expense),
		})
	}
}

func ListExpensesHandler(prv *providers.Provider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := middlewares.GetUserID(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Not authenticated",
			})
		}

		filter := expenseService.ListExpenseFilter{
			Category: c.Query("category"),
			Sort:     c.Query("sort", "date_desc"),
		}

		expenses, err := prv.S.Expense.List(c.Context(), userID, filter)
		if err != nil {
			log.Errorw("Failed to list expenses", "user_id", userID, "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Failed to fetch expenses",
			})
		}

		var totalPaisa int64
		items := make([]fiber.Map, 0, len(expenses))
		for _, e := range expenses {
			totalPaisa += e.Amount
			items = append(items, toExpenseResponse(&e))
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data": fiber.Map{
				"expenses":    items,
				"total":       formatPaisa(totalPaisa),
				"total_paisa": totalPaisa,
				"count":       len(items),
			},
		})
	}
}

func toExpenseResponse(e *models.Expense) fiber.Map {
	return fiber.Map{
		"id":          e.ID,
		"amount":      formatPaisa(e.Amount),
		"amount_paisa": e.Amount,
		"category":    e.Category,
		"description": e.Description,
		"date":        e.Date.Format("2006-01-02"),
		"created_at":  e.CreatedAt,
	}
}

func formatPaisa(paisa int64) string {
	rupees := paisa / 100
	paise := paisa % 100
	if paise < 0 {
		paise = -paise
	}
	return fmt.Sprintf("%d.%02d", rupees, paise)
}