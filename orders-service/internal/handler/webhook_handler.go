package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kyma-project/examples/orders-service/internal/service/model"
	"github.com/kyma-project/examples/orders-service/internal/store"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/kyma-project/examples/orders-service/internal/service"
)

type Webhook struct {
	svc *service.Order
	eventTypes []string
}

func NewWebhook(svc *service.Order) *Webhook {
	return &Webhook{
		svc: svc,
		eventTypes: retrieveEventTypes(),
	}
}

func (h *Webhook) RegisterAll(root string, router Router) {
	router.HandleFunc(fmt.Sprintf("%s", root), h.onHook).Methods(http.MethodPost)
}

func (h *Webhook) onHook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	order := new(model.Order)
	if err := json.Unmarshal(body, order); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.svc.Create(ctx, order)
	if err == store.AlreadyExistsError {
		w.WriteHeader(http.StatusConflict)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Webhook) containsEvent(event string) bool {
	for _, e := range h.eventTypes {
		if e == event {
			return true
		}
	}
	return false
}

func retrieveEventTypes() []string {
	events := os.Getenv("APP_EVENT_TYPES")
	if events == "" {
		return []string{"order.create"}
	}
	return strings.Split(events, ",")
}
