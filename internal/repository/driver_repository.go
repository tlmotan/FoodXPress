package repository

import (
	"context"
	"errors"
	"fmt"
	"grabfood-clone/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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

type PostgresDriverRepository struct {
	db *pgxpool.Pool
}

func NewPostgresDriverRepository(db *pgxpool.Pool) *PostgresDriverRepository {
	// every part of the app shares the exact same pool
	return &PostgresDriverRepository{db: db} // The & symbol means "take the MEMORY ADDRESS of"
	// *PostgresDriverRepository means the the variable that stores the memory address of PostgresDriverRepository
}

func (r *PostgresDriverRepository) GetByID(ctx context.Context, id string) (*domain.Driver, error) {
	query := `
		SELECT id, name, status,
		       ST_Y(location::geometry) AS lat,
		       ST_X(location::geometry) AS lng,
		       updated_at
		FROM drivers
		WHERE id = $1
	`

	var d domain.Driver
	err := r.db.QueryRow(ctx, query, id).Scan(
		&d.ID, &d.Name, &d.Status, &d.Lat, &d.Lng, &d.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// %w keeps the sentinel comparable via errors.Is while still
			// carrying the id for logs.
			return nil, fmt.Errorf("get driver %s: %w", id, domain.ErrDriverNotFound)
		}
		return nil, fmt.Errorf("get driver %s: %w", id, err)
	}
	return &d, nil
}

func (r *PostgresDriverRepository) FindNearestAvailable(ctx context.Context, lat, lng float64, limit int) ([]domain.Driver, error) {
	if limit <= 0 {
		return []domain.Driver{}, nil
	}

	query := `
		SELECT id, name, status,
		       ST_Y(location::geometry) AS lat,
		       ST_X(location::geometry) AS lng,
		       updated_at
		FROM drivers
		WHERE status = $1
		ORDER BY location <-> ST_SetSRID(ST_MakePoint($2, $3), 4326)::geography
		LIMIT $4
	`

	rows, err := r.db.Query(ctx, query, domain.DriverAvailable, lng, lat, limit)
	if err != nil {
		return nil, fmt.Errorf("find nearest available drivers: %w", err)
	}
	defer rows.Close()

	drivers := make([]domain.Driver, 0, limit)
	for rows.Next() {
		var d domain.Driver
		if err := rows.Scan(&d.ID, &d.Name, &d.Status, &d.Lat, &d.Lng, &d.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan nearest available driver: %w", err)
		}
		drivers = append(drivers, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate nearest available drivers: %w", err)
	}

	return drivers, nil
}

func (r *PostgresDriverRepository) UpdateStatus(ctx context.Context, id string, status domain.DriverStatus) error {
	query := `
		UPDATE drivers
		SET status = $2, updated_at = now()
		WHERE id = $1;
	`
	
	tag, err := r.db.Exec(ctx, query, id, status)
	
	if err != nil {
		return fmt.Errorf("update driver %s status: %w", id, err)
	}

	if tag.RowsAffected() == 0{
		return fmt.Errorf("update driver %s status: %w", id, domain.ErrDriverNotFound)
	}

	return nil
}

func (r *PostgresDriverRepository) UpdateLocation(ctx context.Context, id string, lat, lng float64) error {
	// use Redis to change the position of the driver every few seconds
	query := `
		UPDATE drivers
		SET location = ST_SetSRID(ST_MakePoint($2, $3), 4326)::geography,
			updated_at = now()
		WHERE id = $1;
	`
	
	tag, err := r.db.Exec(ctx, query, id, lng, lat)
	
	if err != nil {
		return fmt.Errorf("update driver %s location: %w", id, err)
	}

	if tag.RowsAffected() ==0 {
		return fmt.Errorf("update driver %s location: %w", id, domain.ErrDriverNotFound)
	}

	return nil
}
