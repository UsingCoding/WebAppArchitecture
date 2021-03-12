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

func (repo *menuItemRepository) AddMenuItem(item model.MenuItem) error {
	const insertMenuItemSql = `INSERT INTO menu_item VALUES(?, ?)`

	binaryUUID, _ := item.ID.MarshalBinary()

	_, err := repo.db.Exec(insertMenuItemSql, binaryUUID, item.Name)
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
	ID   uuid.UUID `db:"menu_item_id"`
	Name string    `db:"name"`
}
