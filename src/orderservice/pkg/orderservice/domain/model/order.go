package model

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Order struct {
	ID                 uuid.UUID
	MenuItemIDs        []uuid.UUID
	OrderedAtTimestamp int64
}

type MenuItem struct {
	ID    uuid.UUID
	Name  string
	Price uint
}

var (
	ErrOrderNotFound    = errors.New("order not found")
	ErrMenuItemNotFound = errors.New("menu item not found")
)

type OrderRepository interface {
	GetNextId() uuid.UUID
	FindOrder(id uuid.UUID) (Order, error)
	AddOrder(order Order) error
	RemoveOrder(id uuid.UUID) error
}

type MenuItemRepository interface {
	GetNextId() uuid.UUID
	FindMenuItem(id uuid.UUID) (MenuItem, error)
	FindMenuItems(ids []uuid.UUID) ([]MenuItem, error)
	AddMenuItem(item MenuItem) error
	RemoveMenuItem(id uuid.UUID) error
}
