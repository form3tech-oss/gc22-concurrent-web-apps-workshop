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
	ID     string     `json:"id"`
	Items  []LineItem `json:"items"`
	Status string     `json:"status"`
	Total  float32    `json:"total"`
}

type Orders struct {
	orders map[string]Order
}

func NewOrders() *Orders {
	return &Orders{
		orders: make(map[string]Order),
	}
}

func (n *Orders) Get(o string) Order {
	return n.orders[o]
}

// Upsert creates or updates a new order
func (n *Orders) Upsert(o Order) Order {
	for _, v := range o.Items {
		o.Total += v.Price
	}

	n.orders[o.ID] = o
	return o
}

type LineItem struct {
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	Quantity int8    `json:"quantity"`
}

func NewLineItem(name string, price float32, quantity int8) LineItem {
	return LineItem{name, price, quantity}
}
