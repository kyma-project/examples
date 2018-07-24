package repository

import (
	"fmt"
)

type orderRepositoryMemory struct {
	Orders map[string]Order
}

// NewOrderRepositoryMemory is used to instantiate and return the DB implementation of the OrderRepository.
func NewOrderRepositoryMemory() OrderRepository {
	return &orderRepositoryMemory{Orders: make(map[string]Order)}
}

func (repository *orderRepositoryMemory) InsertOrder(order Order) error {
	id := mapID(order)
	if _, exists := repository.Orders[id]; exists {
		return ErrDuplicateKey
	}
	repository.Orders[id] = order
	return nil
}

func (repository *orderRepositoryMemory) GetOrders() ([]Order, error) {
	ret := make([]Order, 0, len(repository.Orders))
	for _, order := range repository.Orders {
		ret = append(ret, order)
	}
	return ret, nil
}

func (repository *orderRepositoryMemory) GetNamespaceOrders(ns string) ([]Order, error) {
	ret := make([]Order, 0, len(repository.Orders))
	for _, order := range repository.Orders {
		if order.Namespace == ns {
			ret = append(ret, order)
		}
	}
	return ret, nil
}

func (repository *orderRepositoryMemory) DeleteOrders() error {
	repository.Orders = make(map[string]Order)
	return nil
}

func (repository *orderRepositoryMemory) cleanUp() error {
	repository.Orders = make(map[string]Order)
	return nil
}

func (repository *orderRepositoryMemory) DeleteNamespaceOrders(ns string) error {
	for _, order := range repository.Orders {
		if order.Namespace == ns {
			delete(repository.Orders, mapID(order))
		}
	}
	return nil
}

func mapID(o Order) string {
	return fmt.Sprintf("%s-%s", o.OrderId, o.Namespace)
}
