package error

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Message string
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	response := Error{Message: "Not Found"}
	json.NewEncoder(w).Encode(response)
}

func BadRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(400)
	response := Error{Message: "Bad Request"}
	json.NewEncoder(w).Encode(response)
}
