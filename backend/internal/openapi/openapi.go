package openapi

import "net/http"

const DOC_LOCATION = "./docs/openapi.json"

func OpenapiHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, DOC_LOCATION)
}
