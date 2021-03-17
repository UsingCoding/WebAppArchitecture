package query

import "github.com/google/uuid"

type MenuItemQueryService interface {
	GetMenuItem(id uuid.UUID) (MenuItemView, error)
	GetMenuItemsForOrderWithQuantity(id uuid.UUID) ([]MenuItemWithQuantityView, error)
}
