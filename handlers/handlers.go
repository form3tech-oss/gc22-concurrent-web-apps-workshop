package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/form3tech-oss/gc22-concurrent-web-apps-workshop/db"
	"github.com/google/uuid"
)

type Handler struct {
	OrdersDB *db.Orders
}

type Response struct {
	Message string   `json:"message,omitempty"`
	Error   error    `json:"error,omitempty"`
	Order   db.Order `json:"order,omitempty"`
}

// Index is invoked by HTTP GET /
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	// Send an HTTP status & a hardcoded message
	resp := &Response{
		Message: "Welcome to the Digital Ice Cream Van!",
	}
	writeResponse(w, http.StatusOK, resp)
}

// OrderShow is invoked by HTTP POST /orders
func (h *Handler) OrderUpsert(w http.ResponseWriter, r *http.Request) {
	// Initialize an order to unmarshal request body into
	var order db.Order
	// Read the request body
	body, err := readRequestBody(r)
	// Handle any errors & write an error HTTP status & response
	if err != nil {
		resp := &Response{
			Error: fmt.Errorf("invalid order body:%v", err),
		}
		writeResponse(w, http.StatusInternalServerError, resp)
	}
	// Unmarshal response to order var
	// Handle any errors & write an error HTTP status & response
	if err := json.Unmarshal(body, &order); err != nil {
		resp := &Response{
			Error: fmt.Errorf("invalid order body:%v", err),
		}
		writeResponse(w, http.StatusUnprocessableEntity, resp)
	}

	order.ID = uuid.NewString()
	order.Status = db.New.String()
	// Call the repository method corresponding to the operation
	h.OrdersDB.Upsert(order)
	resp := &Response{
		Order: order,
	}
	// Send an HTTP success status & the return value from the repo
	writeResponse(w, http.StatusOK, resp)
}

// writeResponse is a helper method that allows to write and HTTP status & response
func writeResponse(w http.ResponseWriter, status int, resp *Response) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		fmt.Fprintf(w, "error encoding resp %v:%s", resp, err)
	}
}

// readRequestBody is a helper method that allows to read a request body and return any errors
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