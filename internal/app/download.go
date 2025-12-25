package app

import (
	"log"
	"net/http"
)

func (s *Server) DownloadZipHandler(w http.ResponseWriter, r *http.Request) {
	// Defer close on request body
	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			log.Fatalf("could not close request body: %v", closeErr)
		}
	}()
}
