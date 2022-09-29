package db

import "fmt"

type InventoryService struct {
	stock map[string]MenuItem
}

type MenuItem struct {
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func NewInventoryService(s map[string]MenuItem) *InventoryService {
	return &InventoryService{stock: s}
}

func (s *InventoryService) PlaceOrder(items []LineItem) (float64, error) {
	for _, v := range items {
		menuItem, ok := s.stock[v.Name]
		if !ok {
			return 0, fmt.Errorf("item not found %s", v.Name)
		}

		if menuItem.Quantity < v.Quantity {
			return 0, fmt.Errorf("insufficient stock, got %d but wanted %d",
				menuItem.Quantity, v.Quantity)
		}
	}

	var total float64
	for _, v := range items {
		menuItem := s.stock[v.Name]
		menuItem.Quantity -= v.Quantity
		s.stock[v.Name] = menuItem
		total += menuItem.Price * float64(v.Quantity)
	}

	return total, nil
}

func (s *InventoryService) GetStock() []MenuItem {
	var items []MenuItem
	for _, v := range s.stock {
		items = append(items, v)
	}
	return items
}
