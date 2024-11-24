package status

import (
	"encoding/json"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(
		struct {
			Status string `json:"status"`
		}{Status: "Healthy"},
	)
}
