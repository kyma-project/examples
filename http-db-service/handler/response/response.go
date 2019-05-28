package response

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"net/http"
)

type Body struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func WriteCodeAndMessage(code int, msg string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(code)
	response := Body{code, msg}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Error sending response", err)
	}
}
