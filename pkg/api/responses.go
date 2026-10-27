package api

import (
	"encoding/json"
	"net/http"
)

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "error in generating the JSON response", http.StatusInternalServerError)
	}
}
