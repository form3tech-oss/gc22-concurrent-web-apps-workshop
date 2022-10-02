package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "embed"

	"github.com/form3tech-oss/gc22-concurrent-web-apps-workshop/db"
	"github.com/form3tech-oss/gc22-concurrent-web-apps-workshop/handlers"
)

//go:embed stock.json
var stockFile []byte

func main() {
	inventory := importStock()
	s := db.NewInventoryService(inventory)
	o := db.NewOrders(s)
	handler := handlers.NewHandler(o, s)

	router := handlers.ConfigureServer(handler)
	fmt.Println("Listening on localhost:3000...")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func importStock() map[string]db.MenuItem {
	var stock []db.MenuItem

	err := json.Unmarshal(stockFile, &stock)
	if err != nil {
		log.Fatal(err)
	}

	m := make(map[string]db.MenuItem)
	for _, s := range stock {
		m[s.Name] = s
	}

	return m
}
