package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"orderservice/pkg/orderservice/domain/model"
)

func NewMenuItemRepository(db *sqlx.DB) model.MenuItemRepository {
	return &menuItemRepository{db: db}
}

type menuItemRepository struct {
	db *sqlx.DB
}

func (repo *menuItemRepository) GetNextId() uuid.UUID {
	return uuid.New()
}

func (repo *menuItemRepository) FindMenuItem(id uuid.UUID) (model.MenuItem, error) {
	const getMenuItemSql = `SELECT * FROM menu_item WHERE menu_item_id = ?`

	menuItem := sqlxMenuItem{}
	binaryUUID, _ := id.MarshalBinary()

	err := repo.db.Select(&menuItem, getMenuItemSql, binaryUUID)
	if err != nil {
		return model.MenuItem{}, errors.WithStack(err)
	}

	return model.MenuItem{
		ID:   menuItem.ID,
		Name: menuItem.Name,
	}, nil
}

func (repo *menuItemRepository) FindMenuItems(ids []uuid.UUID) ([]model.MenuItem, error) {
	const getMenuItemsSql = `SELECT * FROM menu_item WHERE menu_item_id IN (?)`

	menuItems := make([]sqlxMenuItem, len(ids))
	binaryUUIDS := make([][]byte, len(ids))
	for i, id := range ids {
		binaryUUIDS[i], _ = id.MarshalBinary()
	}

	preparedSql, args, err := sqlx.In(getMenuItemsSql, binaryUUIDS)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = repo.db.Select(&menuItems, preparedSql, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var items []model.MenuItem

	for _, item := range menuItems {
		items = append(items, model.MenuItem{
			ID:   item.ID,
			Name: item.Name,
		})
	}

	return items, nil
}

func (repo *menuItemRepository) AddMenuItem(item model.MenuItem) error {
	const insertMenuItemSql = `INSERT INTO menu_item VALUES(?, ?, ?)`

	binaryUUID, _ := item.ID.MarshalBinary()

	_, err := repo.db.Exec(insertMenuItemSql, binaryUUID, item.Name, item.Price)
	if err != nil {
		return err
	}

	return nil
}

func (repo *menuItemRepository) RemoveMenuItem(id uuid.UUID) error {
	const deleteSql = `DELETE FROM menu_item WHERE menu_item_id = ?`

	binaryUUID, _ := id.MarshalBinary()

	_, err := repo.db.Exec(deleteSql, binaryUUID)
	if err != nil {
		return err
	}

	return nil
}

type sqlxMenuItem struct {
	ID    uuid.UUID `db:"menu_item_id"`
	Name  string    `db:"name"`
	Price uint      `db:"price"`
}
