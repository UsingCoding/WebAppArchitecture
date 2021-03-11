package query

import "github.com/google/uuid"

type OrderQueryService interface {
	GetOrderView(id uuid.UUID) (OrderView, error)
}
