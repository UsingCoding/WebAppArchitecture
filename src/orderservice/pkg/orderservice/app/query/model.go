package query

import "github.com/google/uuid"

type OrderView struct {
	ID                 uuid.UUID
	Items              []MenuItemWithQuantityView
	OrderedAtTimestamp int64
	Cost               uint
}

type MenuItemView struct {
	ID    uuid.UUID
	Name  string
	Price uint
}

type MenuItemWithQuantityView struct {
	MenuItemView
	Quantity uint
}
