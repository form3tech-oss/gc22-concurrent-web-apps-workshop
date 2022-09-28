package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

// ConfigureServer configures the routes of this server and binds handler functions to them
func ConfigureServer(handler *Handler) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.Methods("GET").Path("/").Handler(http.HandlerFunc(handler.Index))
	router.Methods("GET").Path("/orders/{id}").Handler(http.HandlerFunc(handler.OrderByID))
	router.Methods("POST").Path("/orders").Handler(http.HandlerFunc(handler.OrderUpsert))
	
	return router
}
