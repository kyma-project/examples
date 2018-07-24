package repository

import (
	"testing"
	"os"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepositoryMemoryIntegration(t *testing.T) {
	runTestsForRepoType("memory", t)
}

func TestRepositoryDbIntegration(t *testing.T) {
	if os.Getenv("Host") == "" {
		t.Skip("skipping test; DB Config not set")
	}
	runTestsForRepoType("mssql", t)
}

func runTestsForRepoType(repositoryType string, t *testing.T) {
	repo, err := Create(repositoryType)
	require.NoError(t, err)

	newOrder := Order{OrderId: "orderId1", Namespace: "N7", Total: 10}

	t.Run("Create and get Order", func(t *testing.T) {
		//when
		err := repo.InsertOrder(newOrder)
		assert.NoError(t, err)

		//then
		resultOrders, err := repo.GetOrders()
		require.NoError(t, err)
		assert.Len(t, resultOrders, 1)
		assert.Equal(t, resultOrders[0].OrderId, "orderId1")
		assert.Equal(t, resultOrders[0].Namespace, "N7")
		assert.Equal(t, resultOrders[0].Total, float64(10))

		resultOrders, err = repo.GetNamespaceOrders("N7")
		require.NoError(t, err)
		assert.Len(t, resultOrders, 1)
		assert.Equal(t, resultOrders[0].OrderId, "orderId1")
		assert.Equal(t, resultOrders[0].Namespace, "N7")
		assert.Equal(t, resultOrders[0].Total, float64(10))
	})

	t.Run("Return error when order already exists", func(t *testing.T) {
		//when
		err := repo.InsertOrder(newOrder)

		//then
		assert.Equal(t, ErrDuplicateKey, err)
	})

	t.Run("Create orders in different namespaces", func(t *testing.T) {
		//when
		o1 := Order{OrderId: "orderId1", Namespace: "N8", Total: 10}
		o2 := Order{OrderId: "orderId1", Namespace: "N9", Total: 10}
		err := repo.InsertOrder(o1)
		assert.NoError(t, err)
		err = repo.InsertOrder(o2)
		assert.NoError(t, err)

		//then
		resultOrders, err := repo.GetNamespaceOrders("N8")
		require.NoError(t, err)
		assert.Len(t, resultOrders, 1)
		assert.Equal(t, resultOrders[0].OrderId, "orderId1")
		assert.Equal(t, resultOrders[0].Namespace, "N8")
		assert.Equal(t, resultOrders[0].Total, float64(10))


		resultOrders, err = repo.GetNamespaceOrders("N9")
		require.NoError(t, err)
		assert.Len(t, resultOrders, 1)
		assert.Equal(t, resultOrders[0].OrderId, "orderId1")
		assert.Equal(t, resultOrders[0].Namespace, "N9")
		assert.Equal(t, resultOrders[0].Total, float64(10))
	})

	t.Run("Delete Namespace Orders", func(t *testing.T) {
		//when
		err := repo.DeleteNamespaceOrders("N7")
		assert.NoError(t, err)

		// no orders in N7
		resultOrders, err := repo.GetNamespaceOrders("N7")
		assert.NoError(t, err)
		assert.Len(t, resultOrders, 0)

		// all other orders still there
		resultOrders, err = repo.GetOrders()
		assert.NoError(t, err)
		assert.Len(t, resultOrders, 2)
	})

	t.Run("Delete Orders", func(t *testing.T) {
		//when
		err := repo.DeleteOrders()
		assert.NoError(t, err)

		resultOrders, err := repo.GetOrders()
		assert.NoError(t, err)
		assert.Len(t, resultOrders, 0)
	})

	repo.cleanUp()
}