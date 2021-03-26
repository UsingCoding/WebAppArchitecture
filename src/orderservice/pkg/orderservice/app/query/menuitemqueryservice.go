package query

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var (
	ErrMenuItemViewNotFound = errors.New("menu item not found")
)

type MenuItemQueryService interface {
	GetMenuItem(id uuid.UUID) (MenuItemView, error)
	GetMenuItemsForOrderWithQuantity(id uuid.UUID) ([]MenuItemWithQuantityView, error)
}
