package app

import (
	"encoding/json"
	"log"
	"net/http"
)

type UploadForm struct {
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Close request body
	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			log.Fatalf("could not close request body: %v", closeErr)
		}
	}()

	// Parse files from upload
	requestReader, err := r.MultipartReader()
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

	// use io.Copy ?
	// what do we want to return to the frontend?
	// the next logical step would be to flatten any directories and store/return a flattened list of all files
}
