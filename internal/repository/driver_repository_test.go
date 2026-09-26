package repository

import (
	"context"
	"errors"
	"grabfood-clone/internal/domain"
	"math"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Integration tests: these hit the real Postgres from docker-compose.
// Run with:
//   DATABASE_URL=postgres://grabfood:grabfood@localhost:5432/grabfood go test ./internal/repository/ -v

// newTestPool connects to DATABASE_URL, or skips the test if it isn't set
// so `go test ./...` still passes on machines without the DB running.
func newTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL not set; skipping integration test")
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("connect to db: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

// insertTestDriver seeds a driver with known values and deletes it when the
// test ends, so the test never depends on whatever else is in the table.
func insertTestDriver(t *testing.T, pool *pgxpool.Pool, name string, status domain.DriverStatus, lat, lng float64) string {
	t.Helper()
	ctx := context.Background()

	var id string
	err := pool.QueryRow(ctx, `
		INSERT INTO drivers (name, status, location)
		VALUES ($1, $2, ST_SetSRID(ST_MakePoint($3, $4), 4326)::geography)
		RETURNING id
	`, name, status, lng, lat).Scan(&id)
	if err != nil {
		t.Fatalf("insert test driver: %v", err)
	}

	t.Cleanup(func() {
		if _, err := pool.Exec(context.Background(), `DELETE FROM drivers WHERE id = $1`, id); err != nil {
			t.Errorf("cleanup test driver %s: %v", id, err)
		}
	})
	return id
}

func TestGetByID(t *testing.T) {
	ctx := context.Background()
	pool := newTestPool(t)
	repo := NewPostgresDriverRepository(pool)

	// Kuala Lumpur — lat and lng are far apart so a swapped ST_X/ST_Y is obvious.
	const wantLat, wantLng = 3.1390, 101.6869
	existingID := insertTestDriver(t, pool, "Test Driver", domain.DriverAvailable, wantLat, wantLng)

	t.Run("existing driver returns all fields", func(t *testing.T) {
		d, err := repo.GetByID(ctx, existingID)
		if err != nil {
			t.Fatalf("GetByID(%s): unexpected error: %v", existingID, err)
		}

		if d.ID != existingID {
			t.Errorf("ID = %s, want %s", d.ID, existingID)
		}
		if d.Name != "Test Driver" {
			t.Errorf("Name = %q, want %q", d.Name, "Test Driver")
		}
		if d.Status != domain.DriverAvailable {
			t.Errorf("Status = %q, want %q", d.Status, domain.DriverAvailable)
		}
		if !almostEqual(d.Lat, wantLat) || !almostEqual(d.Lng, wantLng) {
			t.Errorf("location = (%f, %f), want (%f, %f)", d.Lat, d.Lng, wantLat, wantLng)
		}
		if d.UpdatedAt.IsZero() {
			t.Error("UpdatedAt is zero, want it populated")
		}
	})

	t.Run("unknown driver returns ErrDriverNotFound", func(t *testing.T) {
		const missingID = "00000000-0000-0000-0000-000000000000"

		d, err := repo.GetByID(ctx, missingID)
		if !errors.Is(err, domain.ErrDriverNotFound) {
			t.Fatalf("GetByID(%s): err = %v, want ErrDriverNotFound", missingID, err)
		}
		if d != nil {
			t.Errorf("GetByID(%s): driver = %+v, want nil", missingID, d)
		}
	})
}

func TestFindNearestAvailable(t *testing.T) {

}

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}
