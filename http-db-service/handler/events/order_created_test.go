package events

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kyma-project/examples/http-db-service/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

const path = "/events/order/created"

func TestHandleOrderCreatedEventSuccess(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc(path, HandleOrderCreatedEvent).Methods(http.MethodPost)

	ts := httptest.NewServer(router)
	defer ts.Close()

	event := repository.OrderCreatedEvent{OrderCode: "handle-success"}
	requestBody := new(bytes.Buffer)
	json.NewEncoder(requestBody).Encode(event)

	res, err := http.Post(ts.URL+path, "application/json", requestBody)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestHandleOrderCreatedEventBadPayload(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc(path, HandleOrderCreatedEvent).Methods(http.MethodPost)

	ts := httptest.NewServer(router)
	defer ts.Close()

	badEvent := struct {
		OrderId string
	}{OrderId: "handle-400"}

	requestBody := new(bytes.Buffer)
	json.NewEncoder(requestBody).Encode(badEvent)

	res, err := http.Post(ts.URL+path, "application/json", requestBody)
	require.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}
