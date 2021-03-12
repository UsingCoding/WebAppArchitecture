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

	err := service.db.Select(&menuItem, getMenuItemSql, binaryUUID)
	if err != nil {
		return query.MenuItemView{}, errors.WithStack(err)
	}

	return query.MenuItemView{
		ID:   menuItem.ID,
		Name: menuItem.Name,
	}, nil
}

type sqlxMenuItem struct {
	ID   uuid.UUID `db:"menu_item_id"`
	Name string    `db:"name"`
}
