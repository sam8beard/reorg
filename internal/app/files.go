package app

import (
	"context"
	"encoding/json"
	//"github.com/minio/minio-go/v7"
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

	// Convert body to int to access upload UUID
	uploadUUID := string(body)
	log.Printf("Should be of type string: %T", uploadUUID)

	// Use upload ID to query all s3_keys of matching rows in files table
	fileRows, dbErr := s.DB.Query(
		context.Background(),
		"SELECT file_name FROM files WHERE upload_uuid = $1",
		uploadUUID,
	)
	if dbErr != nil {
		log.Printf("error from db query call: %v", dbErr)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, dbErr.Error(), http.StatusInternalServerError)
		encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": "could not retrieve files"})
		if encodeErr != nil {
			log.Printf("failed to write response: %v", encodeErr)
			return
		}
		return
	}

	// Scan all rows and store s3 keys
	fileNames := make([]string, 0)
	for fileRows.Next() {
		log.Println("Scanning row...")
		var name string
		if dbErr := fileRows.Scan(&name); dbErr != nil {
			log.Printf("error scanning value from row: %v", dbErr)
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, dbErr.Error(), http.StatusInternalServerError)
			encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": "could not retrieve files"})
			if encodeErr != nil {
				log.Printf("failed to write response: %v", encodeErr)
				return
			}
			return
		}
		fileNames = append(fileNames, name)
	}
	log.Printf("file names: %v", fileNames)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(fileNames); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
