package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Email     *string   `db:"email" json:"email,omitempty"`
	Password  *string   `db:"password" json:"-"`
	Name      string    `db:"name" json:"name"`
	IsGuest   bool      `db:"is_guest" json:"is_guest"`
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}