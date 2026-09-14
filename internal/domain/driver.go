package domain

import "time"

type DriverStatus string

const (
	DriverOffline   DriverStatus = "offline"
	DriverAvailable DriverStatus = "available"
	DriverAssigned  DriverStatus = "assigned"
)

type Driver struct {
	ID        string
	Name      string
	Status    DriverStatus
	Lat       float64
	Lng       float64
	UpdatedAt time.Time
}
