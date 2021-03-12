package model

import "github.com/google/uuid"

type Order struct {
	ID                 uuid.UUID
	Items              []MenuItem
	OrderedAtTimestamp int64
	Cost               int
}

type MenuItem struct {
	ID   uuid.UUID
	Name string
}

type OrderRepository interface {
	GetNextId() uuid.UUID
	FindOrder(id uuid.UUID) (Order, error)
	AddOrder(order Order) error
	RemoveOrder(id uuid.UUID)
}

type MenuItemRepository interface {
	GetNextId() uuid.UUID
	FindMenuItem(id uuid.UUID) (MenuItem, error)
	AddMenuItem(item MenuItem) error
	RemoveMenuItem(id uuid.UUID) error
}
