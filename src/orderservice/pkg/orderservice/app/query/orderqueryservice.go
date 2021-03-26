package query

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

type OrderQueryService interface {
	GetOrderView(id uuid.UUID) (OrderView, error)
	GetOrderViews() ([]OrderView, error)
}
