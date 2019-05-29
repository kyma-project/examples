package handler

import (
	"encoding/json"
	"fmt"
	"github.com/kyma-project/examples/http-db-service/handler/response"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	"github.com/kyma-project/examples/http-db-service/internal/repository"
)

const defaultNamespace = "default"

// Order is used to expose the Order service's basic operations using the HTTP route handler methods which extend it.
type Order struct {
	repository repository.OrderRepository
}

// NewOrderHandler creates a new 'OrderHandler' which provides route handlers for the given OrderRepository's operations.
func NewOrderHandler(repository repository.OrderRepository) Order {
	return Order{repository}
}

// InsertOrder handles an http request for creating an Order given in JSON format.
// The handler also validates the Order payload fields and handles duplicate entry or unexpected errors.
func (orderHandler Order) InsertOrder(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Error parsing request.", err)
		response.WriteCodeAndMessage(http.StatusInternalServerError, "Internal error.", w)
		return
	}

	defer r.Body.Close()
	var order repository.Order
	err = json.Unmarshal(b, &order)
	if err != nil || order.OrderId == "" || order.Total == 0 {
		response.WriteCodeAndMessage(http.StatusBadRequest, "Invalid request body, orderId / total fields cannot be empty.", w)
		return
	}
	if order.Namespace == "" {
		order.Namespace = defaultNamespace
	}

	log.Debugf("Inserting order: '%+v'.", order)
	err = orderHandler.repository.InsertOrder(order)

	switch err {
	case nil:
		w.WriteHeader(http.StatusCreated)
	case repository.ErrDuplicateKey:
		response.WriteCodeAndMessage(http.StatusConflict, fmt.Sprintf("Order %s already exists.", order.OrderId), w)
	default:
		log.Error(fmt.Sprintf("Error inserting order: '%+v'", order), err)
		response.WriteCodeAndMessage(http.StatusInternalServerError, "Internal error.", w)
	}
}

// GetOrders handles an http request for retrieving all Orders from all namespaces.
// The orders list is marshalled in JSON format and sent to the `http.ResponseWriter`
func (orderHandler Order) GetOrders(w http.ResponseWriter, r *http.Request) {
	log.Debug("Retrieving orders")

	orders, err := orderHandler.repository.GetOrders()
	if err != nil {
		log.Error("Error retrieving orders.", err)
		response.WriteCodeAndMessage(http.StatusInternalServerError, "Internal error.", w)
		return
	}

	if err = respondOrders(orders, w); err != nil {
		log.Error("Error sending orders response.", err)
		response.WriteCodeAndMessage(http.StatusInternalServerError, "Internal error.", w)
		return
	}
}

// GetNamespaceOrders handles an http request for retrieving all Orders from a namespace specified as a path variable.
// The orders list is marshalled in JSON format and sent to the `http.ResponseWriter`.
func (orderHandler Order) GetNamespaceOrders(w http.ResponseWriter, r *http.Request) {
	ns, exists := mux.Vars(r)["namespace"]
	if !exists {
		response.WriteCodeAndMessage(http.StatusBadRequest, "No namespace provided.", w)
		return
	}

	log.Debugf("Retrieving orders for namespace: %s\n", ns)

	orders, err := orderHandler.repository.GetNamespaceOrders(ns)
	if err != nil {
		log.Error("Error retrieving orders.", err)
		response.WriteCodeAndMessage(http.StatusInternalServerError, "Internal error.", w)
		return
	}

	if err = respondOrders(orders, w); err != nil {
		log.Error("Error sending orders response.", err)
		response.WriteCodeAndMessage(http.StatusInternalServerError, "Internal error.", w)
		return
	}
}

func respondOrders(orders []repository.Order, w http.ResponseWriter) error {
	body, err := json.Marshal(orders)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(body); err != nil {
		return err
	}
	return nil
}

// DeleteOrders handles an http request for deleting all Orders from all namespaces.
func (orderHandler Order) DeleteOrders(w http.ResponseWriter, r *http.Request) {
	log.Debug("Deleting all orders")

	if err := orderHandler.repository.DeleteOrders(); err != nil {
		log.Error("Error deleting orders.", err)
		response.WriteCodeAndMessage(http.StatusInternalServerError, "Internal error.", w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteNamespaceOrders handles an http request for deleting all Orders from a namespace specified as a path variable.
func (orderHandler Order) DeleteNamespaceOrders(w http.ResponseWriter, r *http.Request) {
	ns, exists := mux.Vars(r)["namespace"]
	if !exists {
		response.WriteCodeAndMessage(http.StatusBadRequest, "No namespace provided.", w)
		return
	}

	log.Debugf("Deleting orders in namespace %s\n", ns)
	if err := orderHandler.repository.DeleteNamespaceOrders(ns); err != nil {
		log.Errorf("Deleting orders in namespace %s\n. %s", ns, err)
		response.WriteCodeAndMessage(http.StatusInternalServerError, "Internal error.", w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
