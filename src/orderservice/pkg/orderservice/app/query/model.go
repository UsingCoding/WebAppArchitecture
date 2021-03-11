package query

import "github.com/google/uuid"

type OrderView struct {
	ID                 uuid.UUID
	Items              []MenuItemView
	OrderedAtTimestamp int64
	Cost               int
}

type MenuItemView struct {
	ID       uuid.UUID
	Quantity int
}
