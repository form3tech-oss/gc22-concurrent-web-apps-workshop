package db

type OrderStatus int

const (
	New OrderStatus = iota
	InProgress
	Completed
)

func (o OrderStatus) String() string {
	return [...]string{"NEW", "IN_PROGRESS", "COMPLETED"}[o]
}

type Order struct {
	ID     string      `json:"id"`
	Items  []string    `json:"items"`
	Status string `json:"status"`
}

type Orders struct {
	orders map[string]Order
}

func NewOrders() *Orders {
	return &Orders{
		orders: make(map[string]Order),
	}
}

// Upsert creates or updates a new order
func (n *Orders) Upsert(o Order) {
	n.orders[o.ID] = o
}
