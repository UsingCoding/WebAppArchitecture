package query

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"orderservice/pkg/orderservice/app/query"
	"strings"
)

func NewOrderQueryService(db *sqlx.DB, menuItemQueryService query.MenuItemQueryService) query.OrderQueryService {
	return &orderQueryService{db: db, menuItemQueryService: menuItemQueryService}
}

type orderQueryService struct {
	db                   *sqlx.DB
	menuItemQueryService query.MenuItemQueryService
}

func (service *orderQueryService) GetOrderView(id uuid.UUID) (query.OrderView, error) {
	orderSql := `SELECT * FROM ||order|| WHERE ||order||.order_id = ?`
	orderSql = strings.Replace(orderSql, "||", "`", -1)
	order := sqlxOrder{}

	binaryUUID, _ := id.MarshalBinary()

	err := service.db.Select(&order, orderSql, binaryUUID)
	if err != nil {
		return query.OrderView{}, errors.WithStack(err)
	}

	menuItems, err := service.menuItemQueryService.GetMenuItemsForOrderWithQuantity(id)
	if err != nil {
		return query.OrderView{}, err
	}

	return query.OrderView{
		ID:                 order.ID,
		Items:              menuItems,
		OrderedAtTimestamp: order.OrderedAtTimestamp,
		Cost:               calculateItemsCost(menuItems),
	}, nil
}

func calculateItemsCost(items []query.MenuItemWithQuantityView) uint {
	var result uint
	for _, item := range items {
		result += item.Price
	}
	return result
}

type sqlxOrder struct {
	ID                 uuid.UUID `db:"id"`
	OrderedAtTimestamp int64     `db:"ordered_at_timestamp"`
}
