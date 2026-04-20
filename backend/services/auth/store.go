
package auth

import (
	"context"

	"github.com/Aman-s12345/expense/backend/db/models"
	"github.com/google/uuid"
)

type Store interface {
	// User operations
	CreateUser(ctx context.Context, user *models.User) error
	FindUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)

	// Session operations
	CreateSession(ctx context.Context, session *models.Session) error
	FindSessionByToken(ctx context.Context, token string) (*models.Session, error)
	DeleteSession(ctx context.Context, id uuid.UUID) error
	DeleteExpiredSessions(ctx context.Context) error
}