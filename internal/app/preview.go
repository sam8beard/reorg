package app

import (
	"log"
	"net/http"
)

func (s *Server) Preview(w http.ResponseWriter, r *http.Request) {
	// Close request body
	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			log.Fatalf("could not close request body: %v", closeErr)
		}
	}()

	// Request body should have json ruleset
	log.Println("receiving request")
}
