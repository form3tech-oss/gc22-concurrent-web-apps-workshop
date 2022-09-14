package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/form3tech-oss/gc22-concurrent-web-apps-workshop/db"
	"github.com/form3tech-oss/gc22-concurrent-web-apps-workshop/handlers"
)

func main() {
	orders := db.NewOrders()
	handler := &handlers.Handler{
		OrdersDB: orders,
	}

	router := handlers.ConfigureServer(handler)
	fmt.Println("Listening on localhost:3000...")
	log.Fatal(http.ListenAndServe(":3000", router))
}
