package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/form3tech-oss/gc22-concurrent-web-apps-workshop/db"
	"github.com/form3tech-oss/gc22-concurrent-web-apps-workshop/handlers"
)

func main() {
	inventory := map[string]int{
		"Solero":     1,
		"Magnum":     1,
		"ScrewBalls": 1,
	}
	s := db.NewInventoryService(inventory)
	handler := &handlers.Handler{
		OrdersDB: db.NewOrders(s),
	}

	router := handlers.ConfigureServer(handler)
	fmt.Println("Listening on localhost:3000...")
	log.Fatal(http.ListenAndServe(":3000", router))
}
