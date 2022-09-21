package db

import "fmt"

type InventoryService struct {
	stock map[string]int
}

func NewInventoryService(s map[string]int) *InventoryService {
	return &InventoryService{stock: s}
}

func (s *InventoryService) DecrementStock(items []LineItem) error {
	for _, v := range items {
		stockLevel, ok := s.stock[v.Name]
		if !ok {
			return fmt.Errorf("item not found %s", v.Name)
		}

		if stockLevel < v.Quantity {
			return fmt.Errorf("insufficent stock, got %d but wanted %d", stockLevel, v.Quantity)
		}
	}

	for _, v := range items {
		stockLevel := s.stock[v.Name]
		stockLevel -= v.Quantity
	}

	return nil
}
