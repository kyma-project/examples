package utils

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"net/http"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func RespondWithCodeAndMessage(code int, msg string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(code)
	response := ErrorResponse{code, msg}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Error sending response", err)
	}
}
