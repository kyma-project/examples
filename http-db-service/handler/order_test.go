package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	responseObj "github.com/kyma-project/examples/http-db-service/handler/response"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"io/ioutil"

	"github.com/kyma-project/examples/http-db-service/internal/repository"
)

func TestCreateOrderSuccess(t *testing.T) {
	// given
	repoMock := repository.MockOrderRepository{}
	defer repoMock.AssertExpectations(t)

	router := mux.NewRouter()
	router.HandleFunc("/orders", NewOrderHandler(&repoMock).InsertOrder).Methods(http.MethodPost)
	ts := httptest.NewServer(router)
	defer ts.Close()

	newOrder := repository.Order{OrderId: "orderId1", Namespace: "N7", Total: 10}

	repoMock.On("InsertOrder", newOrder).Return(nil)

	requestBody := new(bytes.Buffer)
	json.NewEncoder(requestBody).Encode(newOrder)

	// when
	res, err := http.Post(ts.URL+"/orders", "application/json", requestBody)
	require.NoError(t, err)

	// then
	assert.Equal(t, 1, len(repoMock.Calls))
	assert.Equal(t, http.StatusCreated, res.StatusCode)
}

func TestCreateOrderValidation(t *testing.T) {
	// given
	repoMock := repository.MockOrderRepository{}
	defer repoMock.AssertExpectations(t)

	router := mux.NewRouter()
	router.HandleFunc("/orders", NewOrderHandler(&repoMock).InsertOrder).Methods(http.MethodPost)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// when
	res, err := http.Post(ts.URL+"/orders", "application/json", nil)
	require.NoError(t, err)

	// then
	var m responseObj.Body

	b, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	defer res.Body.Close()

	json.Unmarshal(b, &m)

	assert.Equal(t, 400, m.Status)
	assert.Equal(t, "Invalid request body, orderId / total fields cannot be empty.", m.Message)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	assert.Equal(t, 0, len(repoMock.Calls))

	assert.Equal(t, "application/json;charset=UTF-8", res.Header.Get("Content-Type"))
}

func TestCreateOrderConflict(t *testing.T) {
	// given
	repoMock := repository.MockOrderRepository{}
	defer repoMock.AssertExpectations(t)

	router := mux.NewRouter()
	router.HandleFunc("/orders", NewOrderHandler(&repoMock).InsertOrder).Methods(http.MethodPost)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// when
	newOrder := repository.Order{OrderId: "orderId1", Namespace: "N7", Total: 10}

	repoMock.On("InsertOrder", newOrder).Return(repository.ErrDuplicateKey).Once()

	requestBody := new(bytes.Buffer)
	json.NewEncoder(requestBody).Encode(newOrder)

	// when
	response, err := http.Post(ts.URL+"/orders", "application/json", requestBody)
	require.NoError(t, err)

	// then
	var errorResponse responseObj.Body

	responseBody, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	defer response.Body.Close()

	json.Unmarshal(responseBody, &errorResponse)

	assert.Equal(t, 409, errorResponse.Status)
	assert.Equal(t, "Order "+newOrder.OrderId+" already exists.", errorResponse.Message)

	assert.Equal(t, 1, len(repoMock.Calls))
	assert.Equal(t, http.StatusConflict, response.StatusCode)
	assert.Equal(t, "application/json;charset=UTF-8", response.Header.Get("Content-Type"))
}

func TestOrderCreateInternalError(t *testing.T) {
	// given
	repoMock := repository.MockOrderRepository{}
	defer repoMock.AssertExpectations(t)

	router := mux.NewRouter()
	router.HandleFunc("/orders", NewOrderHandler(&repoMock).InsertOrder).Methods(http.MethodPost)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// when
	newOrder := repository.Order{OrderId: "orderId1", Namespace: "N7", Total: 10}

	repoMock.On("InsertOrder", newOrder).Return(errors.New("an error")).Once()

	requestBody := new(bytes.Buffer)
	json.NewEncoder(requestBody).Encode(newOrder)

	response, err := http.Post(ts.URL+"/orders", "application/json", requestBody)
	require.NoError(t, err)

	// then
	responseBody, err := ioutil.ReadAll(response.Body)

	require.NoError(t, err)
	defer response.Body.Close()

	var errorResponse responseObj.Body
	json.Unmarshal(responseBody, &errorResponse)

	assert.Equal(t, 500, errorResponse.Status)
	assert.Equal(t, "Internal error.", errorResponse.Message)

	assert.Equal(t, 1, len(repoMock.Calls))
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	assert.Equal(t, "application/json;charset=UTF-8", response.Header.Get("Content-Type"))
}

func TestCreateOrderWithoutNamespaceSuccess(t *testing.T) {
	// given
	repoMock := repository.MockOrderRepository{}
	defer repoMock.AssertExpectations(t)

	router := mux.NewRouter()
	router.HandleFunc("/orders", NewOrderHandler(&repoMock).InsertOrder).Methods(http.MethodPost)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// insert order without a namespace
	newOrder := repository.Order{OrderId: "orderId1", Total: 10}

	// repository gets an order with default namespace (handler adds it)
	expectedOrder := repository.Order{OrderId: "orderId1", Namespace: "default", Total: 10}
	repoMock.On("InsertOrder", expectedOrder).Return(nil).Once()

	requestBody := new(bytes.Buffer)
	json.NewEncoder(requestBody).Encode(newOrder)

	// when
	res, err := http.Post(ts.URL+"/orders", "application/json", requestBody)
	require.NoError(t, err)

	// then
	assert.Equal(t, 1, len(repoMock.Calls))
	assert.Equal(t, http.StatusCreated, res.StatusCode)
}

func TestGetOrdersSuccess(t *testing.T) {
	// given
	repoMock := repository.MockOrderRepository{}
	defer repoMock.AssertExpectations(t)

	router := mux.NewRouter()
	router.HandleFunc("/orders", NewOrderHandler(&repoMock).GetOrders).Methods(http.MethodGet)
	ts := httptest.NewServer(router)
	defer ts.Close()

	ret := make([]repository.Order, 0)
	repoMock.On("GetOrders").Return(ret, nil).Once()

	// when
	res, err := http.Get(fmt.Sprintf("%s/orders", ts.URL))
	require.NoError(t, err)

	// then
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 1, len(repoMock.Calls))
	assert.Equal(t, "application/json;charset=UTF-8", res.Header.Get("Content-Type"))
}

func TestGetOrdersByNamespaceSuccess(t *testing.T) {
	repoMock := repository.MockOrderRepository{}
	defer repoMock.AssertExpectations(t)

	router := mux.NewRouter()
	router.HandleFunc("/namespace/{namespace}/orders", NewOrderHandler(&repoMock).GetNamespaceOrders).Methods(http.MethodGet)

	ts := httptest.NewServer(router)
	defer ts.Close()
	testNS := "test-namespace"

	// repo mock expects to be passed the namespace in the URL as parameter by the handler, otherwise the test will fail
	ret := make([]repository.Order, 0)
	repoMock.On("GetNamespaceOrders", testNS).Return(ret, nil).Once()

	resp, err := http.Get(fmt.Sprintf("%s/namespace/%s/orders", ts.URL, testNS))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 1, len(repoMock.Calls))
	assert.Equal(t, "application/json;charset=UTF-8", resp.Header.Get("Content-Type"))
}

func TestGetOrderInternalError(t *testing.T) {
	// given
	repoMock := repository.MockOrderRepository{}
	defer repoMock.AssertExpectations(t)

	router := mux.NewRouter()
	router.HandleFunc("/orders", NewOrderHandler(&repoMock).GetOrders).Methods(http.MethodGet)
	ts := httptest.NewServer(router)
	defer ts.Close()

	ret := make([]repository.Order, 0, 0)
	repoMock.On("GetOrders").Return(ret, errors.New("an error")).Once()

	// when
	res, err := http.Get(ts.URL + "/orders")
	require.NoError(t, err)

	// then
	var m responseObj.Body

	b, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	defer res.Body.Close()

	json.Unmarshal(b, &m)

	assert.Equal(t, 500, m.Status)
	assert.Equal(t, "Internal error.", m.Message)

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Equal(t, 1, len(repoMock.Calls))
}

func TestDeletingOrdersSuccess(t *testing.T) {
	repoMock := repository.MockOrderRepository{}
	defer repoMock.AssertExpectations(t)

	router := mux.NewRouter()
	router.HandleFunc("/orders", NewOrderHandler(&repoMock).DeleteOrders).Methods(http.MethodDelete)
	ts := httptest.NewServer(router)
	defer ts.Close()

	repoMock.On("DeleteOrders").Return(nil).Once()

	// when
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/orders", ts.URL), nil)
	res, err := http.DefaultClient.Do(req)

	// then
	require.NoError(t, err)

	assert.Equal(t, http.StatusNoContent, res.StatusCode)
	assert.Equal(t, 1, len(repoMock.Calls))
}

func TestDeleteOrdersByNamespaceSuccess(t *testing.T) {
	repoMock := repository.MockOrderRepository{}
	defer repoMock.AssertExpectations(t)

	router := mux.NewRouter()
	router.HandleFunc("/namespace/{namespace}/orders", NewOrderHandler(&repoMock).DeleteNamespaceOrders).Methods(http.MethodDelete)

	ts := httptest.NewServer(router)
	defer ts.Close()

	testNS := "test-namespace"

	// repo mock expects to be passed the namespace in the URL as parameter by the handler, otherwise the test will fail
	repoMock.On("DeleteNamespaceOrders", testNS).Return(nil).Once()

	// when
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/namespace/%s/orders", ts.URL, testNS), nil)
	// then
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, res.StatusCode)
	assert.Equal(t, 1, len(repoMock.Calls))
}

func TestDeletingOrdersInternalError(t *testing.T) {
	repoMock := repository.MockOrderRepository{}
	defer repoMock.AssertExpectations(t)

	router := mux.NewRouter()
	router.HandleFunc("/orders", NewOrderHandler(&repoMock).DeleteOrders).Methods(http.MethodDelete)
	ts := httptest.NewServer(router)
	defer ts.Close()

	repoMock.On("DeleteOrders").Return(errors.New("an error")).Once()

	// when
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/orders", ts.URL), nil)
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)

	// then
	require.NoError(t, err)
	var m responseObj.Body

	b, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	defer res.Body.Close()

	json.Unmarshal(b, &m)

	assert.Equal(t, 500, m.Status)
	assert.Equal(t, "Internal error.", m.Message)

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Equal(t, 1, len(repoMock.Calls))
}
