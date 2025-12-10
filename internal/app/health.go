package app

import (
	"encoding/json"
	"net/http"
)

/*
Check server health
*/
func (s *Server) HealthHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]bool{"healthy": true}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
