package events

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/kyma-project/examples/http-db-service/handler/response"
	"github.com/kyma-project/examples/http-db-service/internal/repository"
	"io/ioutil"
	"net/http"
)

func HandleOrderCreatedEvent(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Error("error parsing request", err)
		response.WriteCodeAndMessage(http.StatusInternalServerError, "Internal error.", w)
		return
	}

	defer r.Body.Close()

	var event repository.OrderCreatedEvent
	err = json.Unmarshal(b, &event)

	if err != nil || event.OrderCode == "" {
		response.WriteCodeAndMessage(http.StatusBadRequest, "Invalid request body, orderCode cannot be empty.", w)
		return
	}
	fmt.Println("handle order create called.")

	log.Infof("Handling event '%+v' with my custom logic..", event)
	//here add any custom logic such as sending an email, gather insights into purchase.

	w.WriteHeader(http.StatusOK)
}
