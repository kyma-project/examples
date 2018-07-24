package handler

import "net/http"

// SwaggerAPIHandler provides the Swagger API specification file
func SwaggerAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/yaml;charset=UTF-8")
	http.ServeFile(w, r, "api.yaml")
}

// SwaggerAPIRedirectHandler is used to redirect from root (of the service) to the file which contains the service's API
func SwaggerAPIRedirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/api.yaml", http.StatusMovedPermanently)
}
