package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"orderservice/pkg/orderservice/app/query"
	"orderservice/pkg/orderservice/domain/model"
	"strings"
)

type OrderRepository interface {
	model.OrderRepository
	query.OrderQueryService
}

func NewOrderRepository(db *sqlx.DB) OrderRepository {
	return &orderRepository{db: db}
}

type orderRepository struct {
	db *sqlx.DB
}

func (repo *orderRepository) GetOrderView(id uuid.UUID) (query.OrderView, error) {
	domainOrder, err := repo.getOrder(id)
	if err != nil {
		return query.OrderView{}, err
	}

	return convertToOrderView(domainOrder), nil
}

func (repo *orderRepository) GetNextId() uuid.UUID {
	return uuid.New()
}

func (repo *orderRepository) FindOrder(id uuid.UUID) (model.Order, error) {
	return repo.getOrder(id)
}

func (repo *orderRepository) AddOrder(order model.Order) error {
	insertSql := `INSERT INTO ||order|| VALUES (?, ?, ?)`
	insertSql = strings.Replace(insertSql, "||", "`", -1)

	binaryUUID, _ := order.ID.MarshalBinary()

	_, err := repo.db.Exec(insertSql, binaryUUID, order.OrderedAtTimestamp, order.Cost)
	if err != nil {
		return err
	}

	return nil
}

func (repo *orderRepository) RemoveOrder(id uuid.UUID) {
	panic("implement me")
}

type sqlxOrder struct {
	ID                 uuid.UUID `db:"id"`
	OrderedAtTimestamp int64
	Cost               int
}

func (repo orderRepository) getOrder(id uuid.UUID) (model.Order, error) {
	const orderSql = `SELECT * FROM order WHERE order.order_id = ?`
	order := sqlxOrder{}

	binaryUUID, _ := id.MarshalBinary()

	err := repo.db.Select(&order, orderSql, binaryUUID)
	if err != nil {
		return model.Order{}, errors.WithStack(err)
	}

	const menuItemSql = `
		SELECT * 
		FROM menu_item
		LEFT JOIN order_has_menu_item
		WHERE order_has_menu_item.order_id = ?`

	var menuItems []sqlxMenuItem
	err = repo.db.Select(&menuItems, menuItemSql, id)
	if err != nil {
		return model.Order{}, errors.WithStack(err)
	}

	domainOrder := model.Order{
		ID:                 order.ID,
		Items:              nil,
		OrderedAtTimestamp: order.OrderedAtTimestamp,
		Cost:               order.Cost,
	}

	for _, menuItem := range menuItems {
		domainOrder.Items = append(domainOrder.Items, model.MenuItem{
			ID: menuItem.ID,
		})
	}

	return domainOrder, err
}

func convertToOrderView(order model.Order) query.OrderView {
	orderItemViews := make([]query.MenuItemView, len(order.Items))
	for i, item := range order.Items {
		orderItemViews[i] = getMenuItemView(item)
	}

	return query.OrderView{
		ID:                 order.ID,
		Items:              orderItemViews,
		OrderedAtTimestamp: order.OrderedAtTimestamp,
		Cost:               order.Cost,
	}
}

func getMenuItemView(item model.MenuItem) query.MenuItemView {
	return query.MenuItemView{
		ID: item.ID,
	}
}
