package expense

import (
	"context"

	"github.com/Aman-s12345/expense/backend/db/models"
	"github.com/google/uuid"
)

type Store interface {
	Create(ctx context.Context, expense *models.Expense) error
	FindAll(ctx context.Context, userID uuid.UUID, filter ListExpenseFilter) ([]models.Expense, error)
	FindByIdempotencyKey(ctx context.Context, userID uuid.UUID, key string) (*models.Expense, error)
}