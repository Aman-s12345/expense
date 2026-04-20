package expense

import (
	"context"

	"github.com/Aman-s12345/expense/backend/db/models"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, userID uuid.UUID, req CreateExpenseRequest) (*models.Expense, error)
	List(ctx context.Context, userID uuid.UUID, filter ListExpenseFilter) ([]models.Expense, error)
}

type CreateExpenseRequest struct {
	Amount         int64  `json:"amount"`          // in paisa
	Category       string `json:"category"`
	Description    string `json:"description"`
	Date           string `json:"date"`            // "2025-01-15"
	IdempotencyKey string `json:"idempotency_key"` // client-generated UUID
}

type ListExpenseFilter struct {
	Category string // empty = no filter
	Sort     string // "date_desc" or empty (default: date_desc)
}