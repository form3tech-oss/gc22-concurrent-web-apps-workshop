package main

import (
	"bytes"
	"log"
	"net/http"
	"sync"
)

const ordersEndpoint string = "http://localhost:3000/orders"
const indexEndpoint string = "http://localhost:3000/"

const orderCount int = 20
const maxOrderAmount int = 15

// Load test the server
func main() {
	log.Println("Welcome to the Digital Ice Cream Van Order Simulator!")
	if err := checkHealth(); err != nil {
		log.Fatalf("Endpoint %s is not up, please start the server before running simulations.", indexEndpoint)
	}

	log.Printf("We will now simulate %d orders, hold onto your hats!\n", orderCount)

	var wg sync.WaitGroup
	for i := 1; i <= orderCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			createRandomOrder(i)
		}(i)
	}
	wg.Wait()
}

func createRandomOrder(orderNumber int) {
	// rand.Seed(time.Now().UnixNano())
	// amount := rand.Intn(maxOrderAmount) + 1
	// product := products[rand.Intn(len(products))]
	// log.Printf("[simulation-%d]: sending order %+v", orderNumber, item)

	// o := db.Order{
	// 	Items: []db.LineItem{
	// 		db.LineItem{
	// 			Name:     "",
	// 			Quantity: 0,
	// 		},
	// 	},
	// }

	b := []byte(`
		"items":[
			{ "name": "Solero", "quantity": 1 },
			{ "name": "Magnum", "quantity": 1 }
		]
	`)

	req, err := http.NewRequest(http.MethodPost, ordersEndpoint, bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	if _, err = client.Do(req); err != nil {
		log.Fatal(err)
	}

	log.Printf("[simulation-%d]: completed", orderNumber)
}

func checkHealth() error {
	c := http.Client{}
	_, err := c.Get(indexEndpoint)
	return err
}
