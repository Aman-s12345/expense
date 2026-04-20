package expense

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Aman-s12345/expense/backend/db/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) Store {
	return &store{db: db}
}

const expenseSelectFields = `id, user_id, amount, category, description, date, idempotency_key, created_at, updated_at`

func (s *store) Create(ctx context.Context, expense *models.Expense) error {
	query := `
		INSERT INTO expenses (id, user_id, amount, category, description, date, idempotency_key, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (user_id, idempotency_key) DO NOTHING
	`

	result, err := s.db.ExecContext(ctx, query,
		expense.ID,
		expense.UserID,
		expense.Amount,
		expense.Category,
		expense.Description,
		expense.Date,
		expense.IdempotencyKey,
		expense.CreatedAt,
		expense.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert expense: %w", err)
	}

	// If ON CONFLICT hit, rows affected = 0. We already handled this via
	// the idempotency check in the service layer, but this is a DB-level safety net.
	rows, _ := result.RowsAffected()
	if rows == 0 && expense.IdempotencyKey != nil {
		// Fetch the existing one so the caller gets the correct ID back
		existing, err := s.FindByIdempotencyKey(ctx, expense.UserID, *expense.IdempotencyKey)
		if err != nil {
			return fmt.Errorf("failed to fetch existing expense after conflict: %w", err)
		}
		if existing != nil {
			*expense = *existing
		}
	}

	return nil
}

func (s *store) FindAll(ctx context.Context, userID uuid.UUID, filter ListExpenseFilter) ([]models.Expense, error) {
	query := fmt.Sprintf(`SELECT %s FROM expenses WHERE user_id = $1`, expenseSelectFields)
	args := []interface{}{userID}
	argIdx := 2

	if filter.Category != "" {
		query += fmt.Sprintf(` AND category = $%d`, argIdx)
		args = append(args, filter.Category)
		argIdx++
	}

	// Default sort: newest first
	switch filter.Sort {
	case "date_asc":
		query += ` ORDER BY date ASC, created_at ASC`
	default:
		query += ` ORDER BY date DESC, created_at DESC`
	}

	var expenses []models.Expense
	if err := s.db.SelectContext(ctx, &expenses, query, args...); err != nil {
		return nil, fmt.Errorf("failed to list expenses: %w", err)
	}

	if expenses == nil {
		expenses = []models.Expense{}
	}

	return expenses, nil
}

func (s *store) FindByIdempotencyKey(ctx context.Context, userID uuid.UUID, key string) (*models.Expense, error) {
	query := fmt.Sprintf(`SELECT %s FROM expenses WHERE user_id = $1 AND idempotency_key = $2`, expenseSelectFields)

	var expense models.Expense
	err := s.db.GetContext(ctx, &expense, query, userID, key)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find expense by idempotency key: %w", err)
	}

	return &expense, nil
}