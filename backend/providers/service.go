package providers

import (
	"github.com/Aman-s12345/expense/backend/services/auth"
	"github.com/Aman-s12345/expense/backend/services/expense"
)

type Services struct {
	Auth         auth.Service
	AuthStore    auth.Store
	Expense      expense.Service
	ExpenseStore expense.Store
}

func NewServices(p *Provider) *Services {
	sqlxDB := p.DB.GetDB()

	// Auth service
	authStore := auth.NewStore(sqlxDB)
	authService := auth.NewService(authStore)

	// Expense service
	expenseStore := expense.NewStore(sqlxDB)
	expenseService := expense.NewService(expenseStore)

	return &Services{
		Auth:         authService,
		AuthStore:    authStore,
		Expense:      expenseService,
		ExpenseStore: expenseStore,
	}
}