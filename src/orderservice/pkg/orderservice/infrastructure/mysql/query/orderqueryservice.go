package query

import (
	"database/sql"
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
		if err != sql.ErrNoRows {
			return query.OrderView{}, query.ErrOrderNotFound
		}
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

func (service *orderQueryService) GetOrderViews() ([]query.OrderView, error) {
	orderSql := `SELECT * FROM ||order||`
	orderSql = strings.Replace(orderSql, "||", "`", -1)
	var orders []sqlxOrder

	err := service.db.Select(&orders, orderSql)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var orderHasMenuItem []sqlxOrderHasMenuItem

	const menuItemsSql = `
		SELECT *
		FROM order_has_menu_item
		LEFT JOIN menu_item mi on order_has_menu_item.menu_item_id = mi.menu_item_id
	`

	err = service.db.Select(&orderHasMenuItem, menuItemsSql)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	orderIDToOrderHasMenuItemMap := map[uuid.UUID][]sqlxOrderHasMenuItem{}

	for _, item := range orderHasMenuItem {
		items := orderIDToOrderHasMenuItemMap[item.OrderID]
		items = append(items, item)
		orderIDToOrderHasMenuItemMap[item.OrderID] = items
	}

	result := make([]query.OrderView, len(orders))

	for i, order := range orders {
		items := convertToMenuItemWithViewSlice(orderIDToOrderHasMenuItemMap[order.ID])
		result[i] = query.OrderView{
			ID:                 order.ID,
			Items:              items,
			OrderedAtTimestamp: order.OrderedAtTimestamp,
			Cost:               calculateItemsCost(items),
		}
	}

	return result, nil
}

func calculateItemsCost(items []query.MenuItemWithQuantityView) uint {
	var result uint
	for _, item := range items {
		result += item.Price * item.Quantity
	}
	return result
}

func convertToMenuItemWithViewSlice(items []sqlxOrderHasMenuItem) []query.MenuItemWithQuantityView {
	result := make([]query.MenuItemWithQuantityView, len(items))

	for i, item := range items {
		result[i] = query.MenuItemWithQuantityView{
			MenuItemView: query.MenuItemView{
				ID:    item.ID,
				Name:  item.MenuItemName,
				Price: item.MenuItemPrice,
			},
			Quantity: item.Quantity,
		}
	}

	return result
}

type sqlxOrder struct {
	ID                 uuid.UUID `db:"id"`
	OrderedAtTimestamp int64     `db:"ordered_at_timestamp"`
}

type sqlxOrderHasMenuItem struct {
	ID            uuid.UUID `db:"id"`
	OrderID       uuid.UUID `db:"order_id"`
	MenuItemName  string    `db:"name"`
	MenuItemPrice uint      `db:"price"`
	Quantity      uint      `db:"quantity"`
}
