package app

import (
	"context"
	"encoding/json"
	//"github.com/minio/minio-go/v7"
	"github.com/sam8beard/reorg/internal/auth/middleware"
	"io"
	"log"
	"net/http"
)

type FileInfo struct {
	Name   string `json:"fileName"`
	FileID string `json:"fileID"`
}

/* Get files that match a upload id in request */
func (s *Server) FileHandler(w http.ResponseWriter, r *http.Request) {
	// Determine whether request came from a guest or registered user
	//userID := r.Context().Value(middleware.CtxKeyUserID).(string)
	isGuest := r.Context().Value(middleware.CtxKeyGuest).(bool)
	log.Printf("Guest session: %v", isGuest)

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

	// Convert body to int to access upload ID
	uploadID := string(body)

	// Use upload ID to query all s3_keys of matching rows in files table
	fileRows, dbErr := s.DB.Query(
		context.Background(),
		"SELECT name, id FROM files WHERE upload_id = $1",
		uploadID,
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
	files := make([]FileInfo, 0)
	for fileRows.Next() {
		log.Println("Scanning row...")
		var f FileInfo

		if dbErr := fileRows.Scan(&f.Name, &f.FileID); dbErr != nil {
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
		files = append(files, f)
	}
	log.Printf("file names: %v", files)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(files); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
