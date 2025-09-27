package web

import (
	"encoding/json"
	"net/http"

	"github.com/marcelofabianov/fault"
)

func Success(w http.ResponseWriter, r *http.Request, status int, data any) {
	writeJSON(w, status, data)
}

func Error(w http.ResponseWriter, r *http.Request, err error) {
	response := fault.ToResponse(err)
	writeJSON(w, response.StatusCode, response)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
