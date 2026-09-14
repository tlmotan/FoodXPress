package repository

import (
	"context"

	"grabfood-clone/internal/domain"
)

// DriverRepository is a starting point — revisit method set once the
// matching engine's actual query patterns are clear (e.g. nearest-N,
// lock-for-update on assignment).
type DriverRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Driver, error)
	FindNearestAvailable(ctx context.Context, lat, lng float64, limit int) ([]domain.Driver, error)
	UpdateStatus(ctx context.Context, id string, status domain.DriverStatus) error
	UpdateLocation(ctx context.Context, id string, lat, lng float64) error
}
