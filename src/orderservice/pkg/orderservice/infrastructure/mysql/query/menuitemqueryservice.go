package query

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"orderservice/pkg/orderservice/app/query"
)

func NewMenuItemQueryService(db *sqlx.DB) query.MenuItemQueryService {
	return &menuItemQueryService{db: db}
}

type menuItemQueryService struct {
	db *sqlx.DB
}

func (service *menuItemQueryService) GetMenuItem(id uuid.UUID) (query.MenuItemView, error) {
	const getMenuItemSql = `SELECT * FROM menu_item WHERE menu_item_id = ?`

	menuItem := sqlxMenuItem{}
	binaryUUID, _ := id.MarshalBinary()

	err := service.db.Get(&menuItem, getMenuItemSql, binaryUUID)
	if err != nil {
		return query.MenuItemView{}, errors.WithStack(err)
	}

	return query.MenuItemView{
		ID:    menuItem.ID,
		Name:  menuItem.Name,
		Price: menuItem.Price,
	}, nil
}

func (service *menuItemQueryService) GetMenuItemsForOrderWithQuantity(id uuid.UUID) ([]query.MenuItemWithQuantityView, error) {
	const sql = `SELECT * 
		FROM menu_item m
		LEFT JOIN order_has_menu_item order_has ON m.menu_item_id = order_has.menu_item_id
		WHERE order_has.order_id = ?`

	var menuItemsWithQuantity []sqlxMenuItemWithQuantity
	binaryUUID, _ := id.MarshalBinary()

	err := service.db.Select(&menuItemsWithQuantity, sql, binaryUUID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := make([]query.MenuItemWithQuantityView, len(menuItemsWithQuantity))

	for i, item := range menuItemsWithQuantity {
		result[i] = query.MenuItemWithQuantityView{
			MenuItemView: query.MenuItemView{
				ID:    item.ID,
				Name:  item.Name,
				Price: item.Price,
			},
			Quantity: item.Quantity,
		}
	}

	return result, nil
}

type sqlxMenuItem struct {
	ID    uuid.UUID `db:"menu_item_id"`
	Name  string    `db:"name"`
	Price uint      `db:"price"`
}

type sqlxMenuItemWithQuantity struct {
	ID       uuid.UUID `db:"menu_item_id"`
	Name     string    `db:"name"`
	Price    uint      `db:"price"`
	Quantity uint      `db:"quantity"`
}
