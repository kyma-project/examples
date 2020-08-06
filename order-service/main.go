package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/kyma-project/examples/order-service/internal/handler"
	"github.com/kyma-project/examples/order-service/internal/service"
	"github.com/kyma-project/examples/order-service/internal/store"
)

const (
	timeout = 15 * time.Second
)

func main() {
	var storage store.Store
	if os.Getenv("REDIS_HOST") != "" && os.Getenv("REDIS_PORT") != "" {
		redisClient := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
			Password: "",
		})
		storage = store.NewRedis(redisClient)
	} else {
		storage = store.NewMemory()
	}

	ordersSvc := service.NewOrders(storage)

	r := mux.NewRouter()
	r.Use(logRequest)

	order := handler.NewOrder(ordersSvc)
	order.RegisterAll("/order", r)

	webhook := handler.NewWebhook(ordersSvc)
	webhook.RegisterAll("/webhook", r)

	log.Println("List of registered endpoints:")
	err := printEndpoints(r)
	if err != nil {
		log.Fatalf("Cannot print registered routes, because: %v", err)
	}

	srv := http.Server{
		Addr:         ":8080",
		Handler:      cors.AllowAll().Handler(r),
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	log.Println(fmt.Sprintf("Listening on %s", srv.Addr))

	onShutdown(srv, timeout)
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(fmt.Sprintf("[%s] %s", r.Method, r.RequestURI))
		next.ServeHTTP(w, r)
	})
}

func printEndpoints(router *mux.Router) error {
	return router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			return err
		}

		methods, err := route.GetMethods()
		if err != nil {
			return err
		}

		log.Println(fmt.Sprintf("Path: %s, Methods: %s", pathTemplate, strings.Join(methods, ",")))
		return nil
	})
}

func onShutdown(srv http.Server, timeout time.Duration) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("Shutting down, bye bye")
	os.Exit(0)
}
