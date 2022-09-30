package db

import (
	"fmt"

	"github.com/google/uuid"
)

type OrderStatus int

const (
	New OrderStatus = iota
	InProgress
	Completed
	Rejected
)

func (o OrderStatus) String() string {
	return [...]string{"NEW", "IN_PROGRESS", "COMPLETED", "REJECTED"}[o]
}

type Order struct {
	ID     string     `json:"id"`
	Items  []LineItem `json:"items"`
	Status string     `json:"status"`
	Total  string     `json:"total,omitempty"`
}

type Orders struct {
	orders    map[string]Order
	inventory *InventoryService
}

func NewOrders(inventory *InventoryService) *Orders {
	return &Orders{
		orders:    make(map[string]Order),
		inventory: inventory,
	}
}

func (n *Orders) Get(id string) (*Order, error) {
	o, ok := n.orders[id]
	if !ok {
		return nil, fmt.Errorf("no order for id %s found", id)
	}
	return &o, nil
}

// Upsert creates or updates a new order
func (n *Orders) Upsert(o Order) (Order, error) {
	o.ID = uuid.NewString()
	o.Status = New.String()
	total, err := n.inventory.PlaceOrder(o.Items)
	if err != nil {
		o.Status = Rejected.String()
		n.orders[o.ID] = o
		return o, err
	}
	o.Total = fmt.Sprintf("%.2f", total)
	o.Status = Completed.String()

	n.orders[o.ID] = o

	return o, nil
}

type LineItem struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}
