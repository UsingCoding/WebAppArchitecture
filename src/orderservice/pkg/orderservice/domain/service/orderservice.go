package service

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"orderservice/pkg/orderservice/domain/model"
	"time"
)

var (
	ErrIncorrectMenuItemsCount = errors.New("menu items must be more than 0")
	ErrSomeMenuItemsNotFound   = errors.New("some menu items not found")
)

type OrderService interface {
	CreateOrder(menuItemsIDs []uuid.UUID) (model.Order, error)
}

func NewOrderService(orderRepository model.OrderRepository, menuItemRepository model.MenuItemRepository) OrderService {
	return &orderService{orderRepository: orderRepository, menuItemRepository: menuItemRepository}
}

type orderService struct {
	orderRepository    model.OrderRepository
	menuItemRepository model.MenuItemRepository
}

func (service *orderService) CreateOrder(menuItemsIDs []uuid.UUID) (model.Order, error) {
	if len(menuItemsIDs) == 0 {
		return model.Order{}, ErrIncorrectMenuItemsCount
	}

	menuItems, err := service.menuItemRepository.FindMenuItems(menuItemsIDs)
	if err != nil {
		return model.Order{}, err
	}

	if len(menuItems) != len(menuItemsIDs) {
		return model.Order{}, ErrSomeMenuItemsNotFound
	}

	order := model.Order{
		ID:                 service.orderRepository.GetNextId(),
		MenuItemIDs:        menuItemsIDs,
		OrderedAtTimestamp: time.Now().Unix(),
	}

	err = service.orderRepository.AddOrder(order)
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}
