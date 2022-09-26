package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/form3tech-oss/gc22-concurrent-web-apps-workshop/db"
	"github.com/form3tech-oss/gc22-concurrent-web-apps-workshop/handlers"
)

const stockPath = "./cmd/server/stock.json"

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

	dir, _ := os.Getwd()
	log.Print(dir)

	file, err := os.Open(stockPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(b, &stock)
	if err != nil {
		log.Fatal(err)
	}

	m := make(map[string]db.MenuItem)
	for _, s := range stock {
		m[s.Name] = s
	}
	return m
}
