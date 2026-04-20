package expense

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Aman-s12345/expense/backend/db/models"
	"github.com/google/uuid"
)

type service struct {
	store Store
}

func NewService(store Store) Service {
	return &service{
		store: store,
	}
}

func (s *service) Create(ctx context.Context, userID uuid.UUID, req CreateExpenseRequest) (*models.Expense, error) {
	if err := validateCreateRequest(req); err != nil {
		return nil, err
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
	}

	// If idempotency key is provided, check for existing expense first.
	// This prevents duplicate entries when the client retries a request.
	if req.IdempotencyKey != "" {
		existing, err := s.store.FindByIdempotencyKey(ctx, userID, req.IdempotencyKey)
		if err != nil {
			return nil, fmt.Errorf("failed to check idempotency: %w", err)
		}
		if existing != nil {
			return existing, nil // return the already-created expense
		}
	}

	expense := &models.Expense{
		ID:          uuid.New(),
		UserID:      userID,
		Amount:      req.Amount,
		Category:    strings.TrimSpace(req.Category),
		Description: strings.TrimSpace(req.Description),
		Date:        date,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if req.IdempotencyKey != "" {
		expense.IdempotencyKey = &req.IdempotencyKey
	}

	if err := s.store.Create(ctx, expense); err != nil {
		return nil, fmt.Errorf("failed to create expense: %w", err)
	}

	return expense, nil
}

func (s *service) List(ctx context.Context, userID uuid.UUID, filter ListExpenseFilter) ([]models.Expense, error) {
	return s.store.FindAll(ctx, userID, filter)
}

func validateCreateRequest(req CreateExpenseRequest) error {
	if req.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	category := strings.TrimSpace(req.Category)
	if category == "" {
		return fmt.Errorf("category is required")
	}

	if req.Date == "" {
		return fmt.Errorf("date is required")
	}

	return nil
}