package models

import (
	"time"

	"github.com/google/uuid"
)

type Expense struct {
	ID          uuid.UUID `db:"id" json:"id"`
	UserID      uuid.UUID `db:"user_id" json:"user_id"`
	Amount      int64     `db:"amount" json:"amount"`      
	Category    string    `db:"category" json:"category"`
	Description string    `db:"description" json:"description"`
	Date        time.Time `db:"date" json:"date"`
	IdempotencyKey *string `db:"idempotency_key" json:"-"`  
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

func (Expense) TableName() string {
	return "expenses"
}

func (e *Expense) AmountInRupees() float64 {
	return float64(e.Amount) / 100.0
}

type ExpenseResponse struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Amount      string    `json:"amount"`      
	AmountPaisa int64     `json:"amount_paisa"` 
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Date        string    `json:"date"`       
	CreatedAt   time.Time `json:"created_at"`
}