package repository

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)

func TestMemoryCreateAndGet(t *testing.T) {
	repo := NewOrderRepositoryMemory()

	//when
	newOrder := Order{OrderId: "orderId1", Namespace: "N7", Total: 10}

	err := repo.InsertOrder(newOrder)

	if err != nil {
		t.Fatalf("Could not access database. '%s'", err)
	}

	resultOrders, err := repo.GetOrders()

	//then
	require.NoError(t, err)
	assert.Len(t, resultOrders, 1)
	assert.Equal(t, resultOrders[0].OrderId, "orderId1")
	assert.Equal(t, resultOrders[0].Total, float64(10))
}

func TestMemoryErrorOnCreateDuplicate(t *testing.T) {
	Repo := NewOrderRepositoryMemory()

	//when
	newOrder := Order{OrderId: "orderId1", Namespace: "N7", Total: 10}

	err := Repo.InsertOrder(newOrder)
	err = Repo.InsertOrder(newOrder)

	//then
	assert.Equal(t, err, ErrDuplicateKey)
}
func TestMemoryGetByNamespaceSuccess(t *testing.T) {
	repo := NewOrderRepositoryMemory()

	//when
	orderN7 := Order{OrderId: "orderId1", Namespace: "N7", Total: 10}
	orderN8 := Order{OrderId: "orderId1", Namespace: "N8", Total: 10}

	err := repo.InsertOrder(orderN7)
	assert.NoError(t, err)
	err = repo.InsertOrder(orderN8)
	assert.NoError(t, err)

	// in total 2 orders
	resultOrders, err := repo.GetOrders()
	assert.NoError(t, err)
	assert.Len(t, resultOrders, 2)

	// 1 in N7
	resultOrders, err = repo.GetNamespaceOrders("N7")
	assert.NoError(t, err)
	assert.Len(t, resultOrders, 1)

	// 1 in N8
	resultOrders, err = repo.GetNamespaceOrders("N8")
	assert.NoError(t, err)
	assert.Len(t, resultOrders, 1)
}

func TestMemoryDeleteSuccess(t *testing.T) {
	repo := NewOrderRepositoryMemory()

	// ensure there is an order to delete
	err := repo.InsertOrder(Order{OrderId: "orderId1", Namespace: "N7", Total: 10})
	assert.NoError(t, err)

	resultOrders, err := repo.GetOrders()
	assert.NoError(t, err)
	assert.Len(t, resultOrders, 1)

	// delete order and ensure it is gone
	err = repo.DeleteOrders()
	assert.NoError(t, err)

	resultOrders, err = repo.GetOrders()
	assert.NoError(t, err)
	assert.Len(t, resultOrders, 0)
}

func TestMemoryDeleteByNamespaceSuccess(t *testing.T) {
	repo := NewOrderRepositoryMemory()

	//when
	orderN7 := Order{OrderId: "orderId1", Namespace: "N7", Total: 10}
	orderN8 := Order{OrderId: "orderId1", Namespace: "N8", Total: 10}

	err := repo.InsertOrder(orderN7)
	err = repo.InsertOrder(orderN8)

	// in total 2 orders
	resultOrders, err := repo.GetOrders()
	assert.NoError(t, err)
	assert.Len(t, resultOrders, 2)

	// delete in N7
	err = repo.DeleteNamespaceOrders("N7")
	assert.NoError(t, err)

	// No orders in N7
	resultOrders, err = repo.GetNamespaceOrders("N7")
	assert.NoError(t, err)
	assert.Len(t, resultOrders, 0)

	// but still there in N8
	resultOrders, err = repo.GetNamespaceOrders("N8")
	assert.NoError(t, err)
	assert.Len(t, resultOrders, 1)

	// delete in N8
	err = repo.DeleteNamespaceOrders("N8")
	assert.NoError(t, err)

	// No orders in N8
	resultOrders, err = repo.GetNamespaceOrders("N7")
	assert.NoError(t, err)
	assert.Len(t, resultOrders, 0)

	// No orders at all
	resultOrders, err = repo.GetOrders()
	assert.NoError(t, err)
	assert.Len(t, resultOrders, 0)
}
