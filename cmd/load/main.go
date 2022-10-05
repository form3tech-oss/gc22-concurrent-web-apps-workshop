package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/form3tech-oss/gc22-concurrent-web-apps-workshop/db"
	"github.com/form3tech-oss/gc22-concurrent-web-apps-workshop/handlers"
)

const ordersEndpoint string = "http://localhost:3000/orders"
const indexEndpoint string = "http://localhost:3000/"

const orderCount int = 50
const maxOrderAmount int = 3

var products = []string{"Solero", "ScrewBalls", "Magnum"}

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

	j, err := json.Marshal(o)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, ordersEndpoint, bytes.NewBuffer(j))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == http.StatusOK {
		log.Printf("[simulation-%d]: order placed [status: %d]", orderNumber, resp.StatusCode)
	} else {
		var r handlers.Response
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err := json.Unmarshal(body, &r); err == nil {
			log.Printf("[simulation-%d]: failed to place order [status: %d, reason: %s]", orderNumber, resp.StatusCode, r.Error)
		}
	}
}

func checkHealth() error {
	c := http.Client{}
	_, err := c.Get(indexEndpoint)
	return err
}
