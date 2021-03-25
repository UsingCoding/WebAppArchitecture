package transport

import "orderservice/pkg/orderservice/app/query"

func convertToMenuItems(views []query.MenuItemWithQuantityView) []orderItem {
	result := make([]orderItem, len(views))

	for i, view := range views {
		result[i] = convertToMenuItem(view)
	}

	return result
}

func convertToMenuItem(view query.MenuItemWithQuantityView) orderItem {
	return orderItem{
		Id:       view.ID,
		Quantity: view.Quantity,
		Name:     view.Name,
		Price:    view.Price,
	}
}
