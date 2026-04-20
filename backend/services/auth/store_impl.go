package auth

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


const userSelectFields = `id, email, password, name, is_guest, is_active, created_at, updated_at`

func (s *store) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, email, password, name, is_guest, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := s.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.Password,
		user.Name,
		user.IsGuest,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

func (s *store) FindUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := fmt.Sprintf(`SELECT %s FROM users WHERE id = $1 AND is_active = true`, userSelectFields)

	var user models.User
	err := s.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}

	return &user, nil
}

func (s *store) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := fmt.Sprintf(`SELECT %s FROM users WHERE email = $1 AND is_active = true`, userSelectFields)

	var user models.User
	err := s.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return &user, nil
}


func (s *store) CreateSession(ctx context.Context, session *models.Session) error {
	query := `
		INSERT INTO sessions (id, user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := s.db.ExecContext(ctx, query,
		session.ID,
		session.UserID,
		session.Token,
		session.ExpiresAt,
		session.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert session: %w", err)
	}

	return nil
}

func (s *store) FindSessionByToken(ctx context.Context, token string) (*models.Session, error) {
	query := `SELECT id, user_id, token, expires_at, created_at FROM sessions WHERE token = $1`

	var session models.Session
	err := s.db.GetContext(ctx, &session, query, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find session by token: %w", err)
	}

	return &session, nil
}

func (s *store) DeleteSession(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM sessions WHERE id = $1`

	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

func (s *store) DeleteExpiredSessions(ctx context.Context) error {
	query := `DELETE FROM sessions WHERE expires_at < NOW()`

	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", err)
	}

	return nil
}