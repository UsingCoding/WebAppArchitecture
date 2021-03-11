package service

import (
	"github.com/pkg/errors"
	"orderservice/pkg/orderservice/domain/model"
	"time"
)

type OrderService interface {
	CreateOrder(items []model.MenuItem, cost int) (model.Order, error)
}

func NewOrderService(repo model.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

var (
	ErrInvalidPrice            = errors.New("invalid price for order")
	ErrIncorrectMenuItemsCount = errors.New("menu items must more than 0")
)

type orderService struct {
	repo model.OrderRepository
}

func (service *orderService) CreateOrder(items []model.MenuItem, cost int) (model.Order, error) {
	if cost < 0 {
		return model.Order{}, ErrInvalidPrice
	}

	//if len(items) == 0 {
	//	return model.Order{}, ErrIncorrectMenuItemsCount
	//}

	order := model.Order{
		ID:                 service.repo.GetNextId(),
		Items:              items,
		OrderedAtTimestamp: time.Now().Unix(),
		Cost:               cost,
	}

	err := service.repo.AddOrder(order)
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}
