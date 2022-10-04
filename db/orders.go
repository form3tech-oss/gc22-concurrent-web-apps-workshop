package db

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// Order contains the fields of our orders.
type Order struct {
	ID     string     `json:"id"`
	Items  []LineItem `json:"items"`
	Status string     `json:"status"`
	Total  string     `json:"total,omitempty"`
}

// LineItem contains all the different items of an Order.
type LineItem struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

// OrderService is our database type.
type OrderService struct {
	orders    map[string]Order
	inventory *InventoryService
}

// Sales contains total sales and revenue of the ice cream van.
type Sales struct {
	CompletedOrders int    `json:"completed_orders"`
	RejectedOrders  int    `json:"rejected_orders"`
	TotalRevenue    string `json:"total_revenue"`
}

// NewOrders initialises the Orders service given an InventoryService.
func NewOrders(inventory *InventoryService) *OrderService {
	return &OrderService{
		orders:    make(map[string]Order),
		inventory: inventory,
	}
}

// Get returns a given order or errro if none exists.
func (os *OrderService) Get(id string) (*Order, error) {
	o, ok := os.orders[id]
	if !ok {
		return nil, fmt.Errorf("no order for id %s found", id)
	}
	return &o, nil
}

// Upsert creates or updates a new order.
func (os *OrderService) Upsert(o Order) (Order, error) {
	o.ID = uuid.NewString()
	o.Status = New.String()
	total, err := os.inventory.PlaceOrder(o.Items)
	if err != nil {
		o.Status = Rejected.String()
		os.orders[o.ID] = o
		return o, err
	}
	o.Total = fmt.Sprintf("%.2f", total)
	o.Status = Completed.String()

	os.orders[o.ID] = o

	return o, nil
}

// GetSales returns the sales stats of the order service
// This is a costly/long running operation.
func (os *OrderService) GetSales() *Sales {
	getRandomSleep(500)
	revenue := 0.0
	sales := Sales{}
	for _, o := range os.orders {
		if o.Status == Completed.String() {
			sales.CompletedOrders++
			t, _ := strconv.ParseFloat(o.Total, 64)
			revenue += t
		}
		if o.Status == Rejected.String() {
			sales.RejectedOrders++
		}
	}
	sales.TotalRevenue = fmt.Sprintf("%.2f", revenue)
	return &sales
}

// getRandomSleep returns a random sleep up to given max amount.
// This function is used to simulate long running/costly operations.
func getRandomSleep(maxMillis int) {
	rand.Seed(time.Now().Unix())
	n := rand.Intn(maxMillis) + 10

	time.Sleep(time.Duration(n) * time.Millisecond)
}
