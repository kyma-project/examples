package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

const (
	targetHost   = "http://http-db-service:8017"
	ordersPath   = "orders"
	nsOrdersPath = "namespace/%s/orders"
)

type order struct {
	OrderId   string  `json:"orderId"`
	Namespace string  `json:"namespace"`
	Total     float64 `json:"total"`
}

// TestAcceptanceOrders runs an acceptance test on the orders endpoint.
// First it tests getting orders, then adding orders and finally deleting.
func TestAcceptanceOrders(t *testing.T) {
	http.DefaultClient.Timeout = 5 * time.Second

	// guarantee that all is deleted also when the test fails.
	defer delete(t, ordersPath)

	// initially no orders
	os := get(t, ordersPath)
	assert.Len(t, os, 0)

	// insert an order to default Namespace
	post(t, ordersPath, order{OrderId: "66", Total: 9000})

	// check the inserted order
	os = get(t, ordersPath)
	assert.Len(t, os, 1)
	assert.Equal(t, "default", os[0].Namespace)

	// insert an order to custom Namespace
	post(t, ordersPath, order{OrderId: "66", Namespace: "N7", Total: 9000})

	// check get by namespace
	os = get(t, fmt.Sprintf(nsOrdersPath, "N7"))
	assert.Len(t, os, 1)

	// check total orders
	os = get(t, ordersPath)
	assert.Len(t, os, 2)

	// delete custom namespace orders
	delete(t, fmt.Sprintf(nsOrdersPath, "N7"))
	// check no orders in namespace
	os = get(t, fmt.Sprintf(nsOrdersPath, "N7"))
	assert.Len(t, os, 0)
	// but still orders in other namespaces
	os = get(t, ordersPath)
	assert.Len(t, os, 1)

	//delete orders
	delete(t, ordersPath)
	// check no orders
	os = get(t, ordersPath)
	assert.Len(t, os, 0)
}

func get(t *testing.T, path string) []order {
	resp, err := http.Get(fmt.Sprintf("%s/%s", targetHost, path))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	b, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	defer resp.Body.Close()

	var os []order
	require.NoError(t, json.Unmarshal(b, &os))

	return os
}

func post(t *testing.T, path string, data interface{}) {
	order, err := json.Marshal(data)
	require.NoError(t, err)
	body := bytes.NewReader(order)
	resp, err := http.Post(fmt.Sprintf("%s/%s", targetHost, path), "application/json", body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func delete(t *testing.T, path string) {
	r, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", targetHost, path), nil)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(r)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}
