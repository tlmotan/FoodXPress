package domain

import "time"

type OrderStatus string

const (
	OrderPending   OrderStatus = "pending"
	OrderMatched   OrderStatus = "matched"
	OrderPickedUp  OrderStatus = "picked_up"
	OrderDelivered OrderStatus = "delivered"
	OrderCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID         string
	MerchantID string
	DriverID   *string
	Status     OrderStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type StockItem struct {
	ID         string
	MerchantID string
	Name       string
	Quantity   int
}
