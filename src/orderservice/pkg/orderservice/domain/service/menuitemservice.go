package service

import "orderservice/pkg/orderservice/domain/model"

type MenuItemService interface {
	CreateMenuItem(name string) (model.MenuItem, error)
}

func NewMenuItemService(repo model.MenuItemRepository) *menuItemService {
	return &menuItemService{repo: repo}
}

type menuItemService struct {
	repo model.MenuItemRepository
}

func (service *menuItemService) CreateMenuItem(name string) (model.MenuItem, error) {
	menuItem := model.MenuItem{
		ID:   service.repo.GetNextId(),
		Name: name,
	}

	err := service.repo.AddMenuItem(menuItem)

	if err != nil {
		return model.MenuItem{}, err
	}

	return menuItem, nil
}
