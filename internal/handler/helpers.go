package handler

import (
	"encoding/json"
	"net/http"
)

// writeJSON sets Content-Type and writes JSON. Nil body for 204 is allowed.
func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if status == http.StatusNoContent || body == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(body)
}
