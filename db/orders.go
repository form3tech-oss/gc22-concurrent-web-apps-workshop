package db

import "fmt"

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
	Total  string     `json:"total"`
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

func (n *Orders) Get(o string) Order {
	return n.orders[o]
}

// Upsert creates or updates a new order
func (n *Orders) Upsert(o Order) (Order, error) {
	var t float32
	for _, v := range o.Items {
		t += v.Price * float32(v.Quantity)
	}

	o.Total = fmt.Sprintf("%.2f", t)
	o.Status = Completed.String()

	err := n.inventory.DecrementStock(o.Items)
	if err != nil {
		o.Status = Rejected.String()
	}

	n.orders[o.ID] = o

	return o, err
}

type LineItem struct {
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	Quantity int     `json:"quantity"`
}
