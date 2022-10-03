package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/form3tech-oss/gc22-concurrent-web-apps-workshop/db"
)

const ordersEndpoint string = "http://localhost:3000/orders"
const indexEndpoint string = "http://localhost:3000/"

const orderCount int = 20
const maxOrderAmount int = 5

var products = []string{"Solero", "Screwball", "Magnum"}

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
	rand.Seed(time.Now().UnixNano())
	q := rand.Intn(maxOrderAmount) + 1
	n := products[rand.Intn(len(products))]

	o := db.Order{
		Items: []db.LineItem{
			{Name: n, Quantity: q},
		},
	}

	log.Printf("[simulation-%d]: sending order %+v", orderNumber, o.Items)

	json, err := json.Marshal(o)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, ordersEndpoint, bytes.NewBuffer(json))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	if _, err = client.Do(req); err != nil {
		log.Fatal(err)
	}

	log.Printf("[simulation-%d]: order placed", orderNumber)
}

func checkHealth() error {
	c := http.Client{}
	_, err := c.Get(indexEndpoint)
	return err
}
