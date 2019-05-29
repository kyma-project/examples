package repository

import "errors"

// Order contains the details of an order entity.
type Order struct {
	OrderId   string  `json:"orderId"`
	Namespace string  `json:"namespace"`
	Total     float64 `json:"total"`
}

// OrderRepository interface defines the basic operations needed for the order service
//
//go:generate mockery -name OrderRepository -inpkg
type OrderRepository interface {
	InsertOrder(o Order) error
	GetOrders() ([]Order, error)
	GetNamespaceOrders(ns string) ([]Order, error)
	DeleteOrders() error
	DeleteNamespaceOrders(ns string) error
	cleanUp() error
}

// ErrDuplicateKey is thrown when there is an attempt to create an order with an OrderId which already is used.
var ErrDuplicateKey = errors.New("Duplicate key")

type OrderCreatedEvent struct {
	OrderCode string `json:"orderCode"`
}
