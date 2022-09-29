package db

import (
	"fmt"
	"sync"
)

// InventoryService is our database type.
type InventoryService struct {
	lock  sync.Mutex
	stock map[string]MenuItem
}

// MenuItem is the type of items available.
type MenuItem struct {
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

// NewInventoryService creates a new service, given an map of
// items to initialise from.
func NewInventoryService(s map[string]MenuItem) *InventoryService {
	return &InventoryService{stock: s}
}

// PlaceOrder verifies quantities and place an order
// for the given slice of LineItem.
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

// GetStock returns the current stock available.
func (s *InventoryService) GetStock() []MenuItem {
	s.lock.Lock()
	defer s.lock.Unlock()
	var items []MenuItem
	for _, v := range s.stock {
		items = append(items, v)
	}
	return items
}
