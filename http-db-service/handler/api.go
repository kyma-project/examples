package handler

import (
	"fmt"
	"net/http"
)

// SwaggerAPIHandler provides the Swagger API specification file
func SwaggerAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/yaml;charset=UTF-8")
	http.ServeFile(w, r, "api.yaml")
}

// SwaggerAPIRedirectHandler is used to redirect from root (of the service) to the file which contains the service's API
func SwaggerAPIRedirectHandler(w http.ResponseWriter, r *http.Request) {
	protocol := "http"
	if r.TLS != nil {
		protocol = "https"
	}

	rd := fmt.Sprintf("%s://editor.swagger.io/#/?url=%s://%s/api.yaml", protocol, protocol, r.Host)
	http.Redirect(w, r, rd, http.StatusPermanentRedirect)
}
