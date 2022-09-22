package db

import "fmt"

type InventoryService struct {
	stock map[string]MenuItem
}

type MenuItem struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

func NewInventoryService(s map[string]MenuItem) *InventoryService {
	return &InventoryService{stock: s}
}

func (s *InventoryService) PlaceOrder(items []LineItem) error {
	for _, v := range items {
		menuItem, ok := s.stock[v.Name]
		if !ok {
			return fmt.Errorf("item not found %s", v.Name)
		}

		if menuItem.Quantity < v.Quantity {
			return fmt.Errorf("insufficent stock, got %d but wanted %d",
				menuItem.Quantity, v.Quantity)
		}
	}

	for _, v := range items {
		menuItem := s.stock[v.Name]
		menuItem.Quantity -= v.Quantity
		s.stock[v.Name] = menuItem
	}

	return nil
}

func (s *InventoryService) GetStock() []MenuItem {
	var items []MenuItem
	for _, v := range s.stock {
		items = append(items, v)
	}
	return items
}
