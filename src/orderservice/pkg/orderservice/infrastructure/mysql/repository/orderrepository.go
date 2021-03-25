package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"orderservice/pkg/orderservice/domain/model"
	"strings"
)

type OrderRepository interface {
	model.OrderRepository
}

func NewOrderRepository(db *sqlx.DB) OrderRepository {
	return &orderRepository{db: db}
}

type orderRepository struct {
	db *sqlx.DB
}

func (repo *orderRepository) GetNextId() uuid.UUID {
	return uuid.New()
}

func (repo *orderRepository) FindOrder(id uuid.UUID) (model.Order, error) {
	const orderSql = `SELECT * FROM order WHERE order.order_id = ?`
	order := sqlxOrder{}

	binaryUUID, _ := id.MarshalBinary()

	err := repo.db.Select(&order, orderSql, binaryUUID)
	if err != nil {
		return model.Order{}, errors.WithStack(err)
	}

	const menuItemSql = `
		SELECT menu_item.menu_item_id
		FROM menu_item
		LEFT JOIN order_has_menu_item
		WHERE order_has_menu_item.order_id = ?`

	var ids []uuid.UUID
	err = repo.db.Select(&ids, menuItemSql, id)
	if err != nil {
		return model.Order{}, errors.WithStack(err)
	}

	return model.Order{
		ID:                 order.ID,
		MenuItemIDs:        ids,
		OrderedAtTimestamp: order.OrderedAtTimestamp,
	}, err
}

func (repo *orderRepository) AddOrder(order model.Order) error {
	insertSql := `INSERT INTO ||order|| VALUES (?, ?)`
	insertSql = strings.Replace(insertSql, "||", "`", -1)

	binaryUUID, _ := order.ID.MarshalBinary()

	_, err := repo.db.Exec(insertSql, binaryUUID, order.OrderedAtTimestamp)
	if err != nil {
		return err
	}

	return nil
}

func (repo *orderRepository) RemoveOrder(id uuid.UUID) error {
	removeSql := `DELETE FROM ||order|| WHERE order_id = ?`
	removeSql = strings.Replace(removeSql, "||", "`", -1)
	binaryOrderID, _ := id.MarshalBinary()

	_, err := repo.db.Exec(removeSql, binaryOrderID)
	if err != nil {
		return err
	}

	return nil
}

func (repo orderRepository) addItemsToOrder(orderID uuid.UUID, menuItemsIDs []uuid.UUID, quantity uint) error {
	const insertSql = `INSERT INTO order_has_menu_item VALUES (?, ?, ?)`
	binaryOrderID, _ := orderID.MarshalBinary()

	for _, menuItemID := range menuItemsIDs {
		binaryMenuItemID, _ := menuItemID.MarshalBinary()
		_, err := repo.db.Exec(insertSql, uuid.New(), binaryOrderID, binaryMenuItemID, quantity)
		if err != nil {
			return err
		}
	}

	return nil
}

type sqlxOrder struct {
	ID                 uuid.UUID `db:"id"`
	OrderedAtTimestamp int64     `db:"ordered_at_timestamp"`
}
