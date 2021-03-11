package model

import "github.com/google/uuid"

type Order struct {
	ID                 uuid.UUID
	Items              []MenuItem
	OrderedAtTimestamp int64
	Cost               int
}

type MenuItem struct {
	ID       uuid.UUID
	Quantity int
}

type OrderRepository interface {
	GetNextId() uuid.UUID
	FindOrder(id uuid.UUID) (Order, error)
	AddOrder(order Order) error
	RemoveOrder(id uuid.UUID)

	GetMenuItems(ids []uuid.UUID) ([]MenuItem, error)
	AddMenuItem(order Order) error
	RemoveMenuItem(id uuid.UUID)
}
