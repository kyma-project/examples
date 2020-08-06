package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kyma-project/examples/orders-service/internal/service"
)

type Webhook struct {
	svc *service.Order
}

func NewWebhook(svc *service.Order) *Webhook {
	return &Webhook{
		svc: svc,
	}
}

func (h *Webhook) RegisterAll(root string, router Router) {
	router.HandleFunc(fmt.Sprintf("%s", root), h.onHook).Methods(http.MethodPost)
}

func (h *Webhook) onHook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	//TODO: add some events handling...

	w.WriteHeader(http.StatusOK)
}
