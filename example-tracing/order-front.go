package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Order struct {
	OrderCode  string  `json:"orderCode"`
	OrderPrice float64 `json:"orderPrice"`
}

type StoredOrder struct {
	OrderId   string  `json:"orderId"`
	Namespace string  `json:"namespace"`
	Total     float64 `json:"total"`
}

var httpTransport *http.Transport
var httpClient *http.Client

func main() {

	var (
		port  = flag.Int("port", 8080, "tcp port on which to listen for http requests")
		dbUrl = flag.String("db-url", "", "db url to which store order data")
	)
	flag.Parse()

	httpTransport = &http.Transport{}
	httpClient = &http.Client{
		Transport: httpTransport,
	}

	http.Handle("/orders", ordersHandler(dbUrl))

	log.Printf("HTTP server starting on port %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

func ordersHandler(dbUrl *string) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch requestMethod := &r.Method; *requestMethod {
		case http.MethodPost:
			if r.Body == nil {
				http.Error(w, "Please send a CloudEvent in the HTTP request body", http.StatusBadRequest)
				return
			}

			b, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				log.Printf("Error reading HTTP request body: %v", err)
				http.Error(w, "Error reading HTTP request body", http.StatusBadRequest)
				return
			}

			var order Order
			if err := json.Unmarshal(b, &order); err != nil {
				log.Printf("Error unmarshalling event data: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			log.Println(r.Header)

			if err := storeOrdersInDB(&order, dbUrl, r.Header); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusCreated)
			}
		case http.MethodGet:
			orders, err := getOrdersFromDB(dbUrl, r.Header)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				json.NewEncoder(w).Encode(orders)
			}
		case http.MethodDelete:
			statusCode, err := deleteOrdersFromDB(dbUrl, r.Header)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.WriteHeader(statusCode)
			}
		default:
			http.Error(w, fmt.Sprintf("HTTP method '%v' is not supported", *requestMethod), http.StatusMethodNotAllowed)
		}
	})
}

func getOrdersFromDB(dbUrl *string, incomingHeaders http.Header) (*[]StoredOrder, error) {
	downstreamRequest, err := http.NewRequest(http.MethodGet, *dbUrl, nil)
	propagateTracingHeaders(incomingHeaders, downstreamRequest)
	resp, err := httpClient.Do(downstreamRequest)
	if err != nil {
		return nil, err
	}

	orders := make([]StoredOrder, 0)
	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(byteArray, &orders)

	return &orders, nil
}

func deleteOrdersFromDB(dbUrl *string, incomingHeaders http.Header) (int, error) {
	downstreamRequest, err := http.NewRequest(http.MethodDelete, *dbUrl, nil)
	propagateTracingHeaders(incomingHeaders, downstreamRequest)
	resp, err := httpClient.Do(downstreamRequest)
	if err != nil {
		return -1, err
	}

	return resp.StatusCode, nil
}

func storeOrdersInDB(order *Order, dbUrl *string, incomingHeaders http.Header) error {
	toSend := StoredOrder{OrderId: order.OrderCode, Total: order.OrderPrice}
	payload, err := json.Marshal(toSend)

	if err != nil {
		return err
	}
	downstreamRequest, err := http.NewRequest(http.MethodPost, *dbUrl, bytes.NewBuffer(payload))
	downstreamRequest.Header.Add("Content-Type", "application/json")
	propagateTracingHeaders(incomingHeaders, downstreamRequest)

	resp, err := httpClient.Do(downstreamRequest)
	statusCode := resp.StatusCode

	if statusCode >= 399 && statusCode != 409 {
		return errors.New("error status when storing event")
	}

	return err
}

func propagateTracingHeaders(incomingHeaders http.Header, downstreamRequest *http.Request) {
	traceHeadersName := [...]string{"X-Request-Id", "X-B3-Traceid", "X-B3-Spanid", "X-B3-Parentspanid", "X-B3-Sampled", "X-B3-Flags", "X-Ot-Span-Context"}
	for _, headerName := range traceHeadersName {
		headerVal := incomingHeaders[headerName]
		if headerVal != nil && len(headerVal) > 0 {
			log.Print(headerName, headerVal)
			downstreamRequest.Header[headerName] = headerVal
		}
	}
}
