package repository

import (
	"context"

	"grabfood-clone/internal/domain"
)

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	GetByID(ctx context.Context, id string) (*domain.Order, error)
	UpdateStatus(ctx context.Context, id string, status domain.OrderStatus) error
	AssignDriver(ctx context.Context, orderID, driverID string) error
}

// StockRepository is where the overselling-prevention logic will live —
// e.g. DecrementIfAvailable should be the single point that enforces
// quantity >= 0 under concurrent access (SELECT FOR UPDATE or similar).
type StockRepository interface {
	DecrementIfAvailable(ctx context.Context, itemID string, qty int) (bool, error)
}
