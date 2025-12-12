package app

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

/* Get files that match a upload id in request */
func (s *Server) FileHandler(w http.ResponseWriter, r *http.Request) {
	// Close request body
	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			log.Fatalf("could not close request body: %v", closeErr)
		}
	}()

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, err.Error(), http.StatusBadRequest)
		encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": "could not read request body"})
		if encodeErr != nil {
			log.Printf("failed to write response: %v", encodeErr)
			return
		}
		return
	}

	// Do I need to unmarshal the body? Or since its a single value can I use it directly?
	log.Printf("%v", body)
}
