package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/Aman-s12345/expense/backend/db/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	sessionDuration     = 7 * 24 * time.Hour
	guestSessionDuration = 30 * 24 * time.Hour 
	bcryptCost          = 12
)

type service struct {
	store Store
}

func NewService(store Store) Service {
	return &service{
		store: store,
	}
}

func (s *service) Register(ctx context.Context, req RegisterRequest) (*models.User, *models.Session, error) {
	if err := validateRegisterRequest(req); err != nil {
		return nil, nil, err
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

	existing, err := s.store.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing != nil {
		return nil, nil, fmt.Errorf("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcryptCost)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to hash password: %w", err)
	}

	hashedStr := string(hashedPassword)
	user := &models.User{
		ID:        uuid.New(),
		Email:     &email,
		Password:  &hashedStr,
		Name:      strings.TrimSpace(req.Name),
		IsGuest:   false,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.store.CreateUser(ctx, user); err != nil {
		return nil, nil, fmt.Errorf("failed to create user: %w", err)
	}

	session, err := s.createSession(ctx, user.ID, sessionDuration)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create session: %w", err)
	}

	return user, session, nil
}

func (s *service) Login(ctx context.Context, req LoginRequest) (*models.User, *models.Session, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" || req.Password == "" {
		return nil, nil, fmt.Errorf("email and password are required")
	}

	user, err := s.store.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, nil, fmt.Errorf("invalid email or password")
	}

	if user.Password == nil {
		return nil, nil, fmt.Errorf("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(req.Password)); err != nil {
		return nil, nil, fmt.Errorf("invalid email or password")
	}

	session, err := s.createSession(ctx, user.ID, sessionDuration)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create session: %w", err)
	}

	return user, session, nil
}

func (s *service) GuestLogin(ctx context.Context) (*models.User, *models.Session, error) {
	user := &models.User{
		ID:        uuid.New(),
		Name:      "Guest",
		IsGuest:   true,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.store.CreateUser(ctx, user); err != nil {
		return nil, nil, fmt.Errorf("failed to create guest user: %w", err)
	}

	session, err := s.createSession(ctx, user.ID, guestSessionDuration)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create session: %w", err)
	}

	return user, session, nil
}

func (s *service) ValidateSession(ctx context.Context, token string) (*models.User, error) {
	if token == "" {
		return nil, fmt.Errorf("token is required")
	}

	session, err := s.store.FindSessionByToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to find session: %w", err)
	}
	if session == nil {
		return nil, fmt.Errorf("invalid or expired session")
	}

	if time.Now().After(session.ExpiresAt) {
		_ = s.store.DeleteSession(ctx, session.ID)
		return nil, fmt.Errorf("session expired")
	}

	user, err := s.store.FindUserByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil || !user.IsActive {
		return nil, fmt.Errorf("user not found or inactive")
	}

	return user, nil
}

func (s *service) Logout(ctx context.Context, token string) error {
	session, err := s.store.FindSessionByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("failed to find session: %w", err)
	}
	if session == nil {
		return nil 
	}

	return s.store.DeleteSession(ctx, session.ID)
}

func (s *service) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.store.FindUserByID(ctx, id)
}

func (s *service) createSession(ctx context.Context, userID uuid.UUID, duration time.Duration) (*models.Session, error) {
	token, err := generateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	session := &models.Session{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(duration),
		CreatedAt: time.Now(),
	}

	if err := s.store.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func validateRegisterRequest(req RegisterRequest) error {
	email := strings.TrimSpace(req.Email)
	if email == "" {
		return fmt.Errorf("email is required")
	}
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return fmt.Errorf("invalid email format")
	}
	if len(req.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}