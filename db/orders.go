package db

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
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
	lock           sync.Mutex
	orders         map[string]Order
	inventory      *InventoryService
	incomingOrders chan Order
	done           chan struct{}
	isClosed       bool
}

// Sales contains total sales and revenue of the ice cream van.
type Sales struct {
	CompletedOrders int    `json:"completed_orders"`
	RejectedOrders  int    `json:"rejected_orders"`
	TotalRevenue    string `json:"total_revenue"`
}

// NewOrders initialises the Orders service given an InventoryService.
func NewOrders(workerCount int, inventory *InventoryService) *OrderService {
	n := &OrderService{
		orders:         make(map[string]Order),
		inventory:      inventory,
		incomingOrders: make(chan Order, workerCount),
		done:           make(chan struct{}),
	}

	for i := 0; i < workerCount; i++ {
		go n.processOrderWorker(fmt.Sprintf("Worker-%d", i), n.incomingOrders, n.done)
	}
	return n
}

func (os *OrderService) processOrderWorker(id string, in <-chan Order,
	done <-chan struct{}) {
	log.Printf("%s started up.", id)
	for {
		select {
		case o := <-in:
			log.Println(o)
			err := os.upsert(o)
			if err != nil {
				log.Println(err)
			}
		case <-done:
			log.Printf("%s shut down.", id)
			return
		}
	}
}

// Get returns a given order or error if none exists.
func (os *OrderService) PlaceOrder(o Order) Order {
	o.ID = uuid.NewString()
	if !os.isClosed {
		o.Status = New.String()
		os.incomingOrders <- o
		return o
	}
	o.Status = Rejected.String()
	return o
}

// Get returns a given order or error if none exists.
func (os *OrderService) Get(id string) (*Order, error) {
	os.lock.Lock()
	defer os.lock.Unlock()
	o, ok := os.orders[id]
	if !ok {
		return nil, fmt.Errorf("no order for id %s found", id)
	}
	return &o, nil
}

// Upsert creates or updates a new order.
func (os *OrderService) upsert(o Order) error {
	os.lock.Lock()
	defer os.lock.Unlock()
	total, err := os.inventory.PlaceOrder(o.Items)
	if err != nil {
		o.Status = Rejected.String()
		os.orders[o.ID] = o
		return err
	}
	o.Total = fmt.Sprintf("%.2f", total)
	o.Status = Completed.String()

	os.orders[o.ID] = o

	return nil
}

// Close closes the orders app for taking any new orders
func (os *OrderService) Close() {
	if !os.isClosed {
		close(os.done)
		os.isClosed = true
	}
}

// GetSales returns the sales stats of the order service
// This is a costly/long running operation.
func (os *OrderService) GetSales() *Sales {
	os.lock.Lock()
	defer os.lock.Unlock()
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
