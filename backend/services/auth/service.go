package auth

import (
	"context"

	"github.com/Aman-s12345/expense/backend/db/models"
	"github.com/google/uuid"
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*models.User, *models.Session, error)
	Login(ctx context.Context, req LoginRequest) (*models.User, *models.Session, error)
	GuestLogin(ctx context.Context) (*models.User, *models.Session, error)
	ValidateSession(ctx context.Context, token string) (*models.User, error)
	Logout(ctx context.Context, token string) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}