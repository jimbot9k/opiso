package error

import (
	"encoding/json"
	"net/http"
)

type error struct {
	Message string
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	response := error{Message: "Not Found"}
	json.NewEncoder(w).Encode(response)
}

func BadRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(400)
	response := error{Message: "Bad Request"}
	json.NewEncoder(w).Encode(response)
}
