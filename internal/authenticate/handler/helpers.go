package handler

import (
	"encoding/json"
	"net/http"

	"github.com/NjiruClinton/reproxy-mobile-bff/internal/types"
)

func writeOK(w http.ResponseWriter, status int, data interface{}) {
	writeJSON(w, status, types.APIResponse{OK: true, Data: data})
}

func writeError(w http.ResponseWriter, status int, msg, code string) {
	writeJSON(w, status, types.APIResponse{OK: false, Error: msg, Code: code})
}

func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
