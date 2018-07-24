package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/vrischmann/envconfig"

	"github.com/kyma-project/examples/http-db-service/config"
	"github.com/kyma-project/examples/http-db-service/handler"
	"github.com/kyma-project/examples/http-db-service/internal/repository"

	_ "github.com/lib/pq"
)

func main() {
	log.Println("Starting service...")

	var cfg config.Service
	if err := envconfig.Init(&cfg); err != nil {
		log.Panicf("Error loading main configuration %v\n", err.Error())
	}
	log.Print(cfg)

	repo, err := repository.Create(cfg.DbType)
	if err != nil {
		log.Fatal("Unable to initiate repository", err)
	}

	orderHandler := handler.NewOrderHandler(repo)

	if err := startService(orderHandler, cfg.Port); err != nil {
		log.Fatal("Unable to start server", err)
	}
}

func startService(orderHandler handler.Order, port string) error {
	router := mux.NewRouter().StrictSlash(true)

	// orders
	router.HandleFunc("/orders", orderHandler.InsertOrder).Methods(http.MethodPost)

	router.HandleFunc("/orders", orderHandler.GetOrders).Methods(http.MethodGet)
	router.HandleFunc("/namespace/{namespace}/orders", orderHandler.GetNamespaceOrders).Methods(http.MethodGet)

	router.HandleFunc("/orders", orderHandler.DeleteOrders).Methods(http.MethodDelete)
	router.HandleFunc("/namespace/{namespace}/orders", orderHandler.DeleteNamespaceOrders).Methods(http.MethodDelete)

	// API
	router.HandleFunc("/", handler.SwaggerAPIRedirectHandler).Methods(http.MethodGet)
	router.HandleFunc("/api.yaml", handler.SwaggerAPIHandler).Methods(http.MethodGet)

	log.Printf("Starting server on port %s ", port)

	c := cors.AllowAll()
	return http.ListenAndServe(":"+port, c.Handler(router))
}
