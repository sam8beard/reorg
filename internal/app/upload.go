package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"io"
	"log"
	"net/http"
)

func (s *Server) UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Generate new upload ID
	uploadUUID := uuid.New()

	// Close request body
	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			log.Fatalf("could not close request body: %v", closeErr)
		}
	}()

	// Create multi part reader to manually stream chunks of files
	multiReader, err := r.MultipartReader()
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

	// Read files in chunks
	for {
		part, err := multiReader.NextPart()
		//if part.FormName() == "user" {
		//	log.Printf("should be user ID: %v", part)
		//}
		if err == io.EOF {
			break
		}
		defer part.Close()
		if part.FileName() != "" && part.FileName() != ".DS_Store" {
			// File is found
			// Insert each file into minio bucket under upload id
			var opts minio.PutObjectOptions
			objKey := fmt.Sprintf("%s/%s", uploadUUID, part.FileName())
			_, minioErr := s.Minio.PutObject(context.Background(), s.MinioBucket, objKey, part, -1, opts)
			if minioErr != nil {
				log.Printf("error from minio put object call: %v", minioErr)
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, minioErr.Error(), http.StatusInternalServerError)
				encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": "could not register upload in object storage"})
				if encodeErr != nil {
					log.Printf("failed to write response: %v", encodeErr)
					return
				}
				return
			}

		} else {
			// Text field found
			data, _ := io.ReadAll(part)
			log.Printf("field: %s \t value: %s\n", part.FormName(), string(data))
		}

	}

	// Insert row in uploads table
	_, dbErr := s.DB.Exec(
		context.Background(),
		"INSERT INTO uploads (upload_uuid, user_id) VALUES ($1, NULL)",
		uploadUUID,
	)
	if dbErr != nil {
		log.Printf("error from db exec call: %v", dbErr)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": "could not register upload with user"})
		if encodeErr != nil {
			log.Printf("failed to write response: %v", encodeErr)
			return
		}
		return
	}

	// Return upload id in response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(uploadUUID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
