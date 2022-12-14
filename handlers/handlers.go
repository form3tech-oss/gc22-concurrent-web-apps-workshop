package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/form3tech-oss/gc22-concurrent-web-apps-workshop/db"
	"github.com/gorilla/mux"
)

// Handler contains the handler and all its dependencies.
type Handler struct {
	Orders    *db.OrderService
	Inventory *db.InventoryService
}

// MenuItem is the externally set menu item type.
type MenuItem struct {
	Name     string `json:"name"`
	Quantity string `json:"quantity"`
}

// Response contains all the response types of our handlers.
type Response struct {
	Message string        `json:"message,omitempty"`
	Menu    []db.MenuItem `json:"menu,omitempty"`
	Error   string        `json:"error,omitempty"`
	Order   *db.Order     `json:"order,omitempty"`
	Sales   *db.Sales     `json:"sales,omitempty"`
}

// NewHandler initialises a new handler, given dependencies.
func NewHandler(o *db.OrderService, i *db.InventoryService) *Handler {
	return &Handler{
		Orders:    o,
		Inventory: i,
	}
}

// Index is invoked by HTTP GET /.
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	// Send an HTTP status & a hardcoded message
	resp := &Response{
		Message: "Welcome to the Digital Ice Cream Van!",
		Menu:    h.Inventory.GetStock(),
	}
	writeResponse(w, http.StatusOK, resp)
}

// OrderByID gets the order by ID provided
func (h *Handler) OrderByID(w http.ResponseWriter, r *http.Request) {
	orderID := mux.Vars(r)["id"]
	order, err := h.Orders.Get(orderID)
	if err != nil {
		writeResponse(w, http.StatusNotFound, &Response{
			Error: err.Error(),
		})
		return
	}

	// Send an HTTP success status & the return value from the repo
	writeResponse(w, http.StatusOK, &Response{
		Order: order,
	})
}

// OrderShow is invoked by HTTP POST /orders.
func (h *Handler) OrderUpsert(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := readRequestBody(r)
	// Handle any errors & write an error HTTP status & response
	if err != nil {
		resp := &Response{
			Error: fmt.Errorf("invalid order body:%v", err).Error(),
		}
		writeResponse(w, http.StatusInternalServerError, resp)
	}

	// Initialize an order to unmarshal request body into
	var order db.Order
	// Unmarshal response to order var
	// Handle any errors & write an error HTTP status & response
	if err := json.Unmarshal(body, &order); err != nil {
		resp := &Response{
			Error: fmt.Errorf("invalid order body:%v", err).Error(),
		}
		writeResponse(w, http.StatusUnprocessableEntity, resp)
	}

	// Call the repository method corresponding to the operation
	order, err = h.Orders.Upsert(order)
	resp := &Response{
		Order: &order,
	}
	if err != nil {
		resp.Error = err.Error()
		writeResponse(w, http.StatusBadRequest, resp)
		return
	}

	// Send an HTTP success status & the return value from the repo
	writeResponse(w, http.StatusOK, resp)
}

// Sales is invoked by GET /sales.
func (h *Handler) Sales(w http.ResponseWriter, r *http.Request) {
	resp := &Response{
		Sales: h.Orders.GetSales(),
	}
	writeResponse(w, http.StatusOK, resp)
}

// writeResponse is a helper method that allows to write and HTTP status & response
func writeResponse(w http.ResponseWriter, status int, resp *Response) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if status != http.StatusOK {
		w.WriteHeader(status)
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		fmt.Fprintf(w, "error encoding resp %v:%s", resp, err)
	}
}

// readRequestBody is a helper method that
// allows to read a request body and return any errors.
func readRequestBody(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return []byte{}, err
	}
	if err := r.Body.Close(); err != nil {
		return []byte{}, err
	}
	return body, err
}
