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
	"strconv"
	"time"
)

type FileMeta struct {
	FileName    string `json:"name"`
	DateCreated string `json:"dateCreated"`
	MimeType    string `json:"mimetype"`
}

func (s *Server) UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Generate new upload UUID
	uploadUUID := uuid.New()

	// Close request body
	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			log.Fatalf("could not close request body: %v", closeErr)
		}
	}()

	// Create multi filePart reader to stream files
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

	// Insert row in uploads table
	_, dbErr := s.DB.Exec(
		context.Background(),
		"INSERT INTO uploads (upload_uuid, user_id) VALUES ($1, NULL)",
		uploadUUID,
	)
	if dbErr != nil {
		log.Printf("error from db exec call: %v", dbErr)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, dbErr.Error(), http.StatusInternalServerError)
		encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": "could not register upload with user"})
		if encodeErr != nil {
			log.Printf("failed to write response: %v", encodeErr)
			return
		}
		return
	}

	// Read files in chunks
	for {
		// Get part
		filePart, err := multiReader.NextPart()
		if err == io.EOF {
			break
		}
		defer filePart.Close()

		// filePart.FormName() = fileName, filePart.FileName() = dateCreated
		if filePart.FormName() != "" && filePart.FormName() != ".DS_Store" {
			// File is found
			// Insert each file into minio bucket under upload id
			var opts minio.PutObjectOptions
			// Generate file uuid for db write
			fileUUID := uuid.New()
			objKey := fmt.Sprintf("%s/%s_%s", uploadUUID, fileUUID, filePart.FormName())
			uploadInfo, minioErr := s.Minio.PutObject(context.Background(), s.MinioBucket, objKey, filePart, -1, opts)
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

			// Get file time creation data
			dateMilliInt, err := strconv.Atoi(filePart.FileName())
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, dbErr.Error(), http.StatusInternalServerError)
				encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": "could not retrieve file creation date"})
				if encodeErr != nil {
					log.Printf("failed to write response: %v", encodeErr)
					return
				}
				return
			}
			dateMilliInt64 := int64(dateMilliInt)
			dateCreated := time.UnixMilli(dateMilliInt64).UTC()
			// Get file size
			fileSize := uploadInfo.Size
			// Get mime type
			mimeType := filePart.Header.Get("Content-Type")

			// Insert row in files table
			_, dbErr := s.DB.Exec(
				context.Background(),
				"INSERT INTO files (upload_id, file_uuid, upload_uuid, name, s3_key, size, mime_type, original_timestamp) VALUES ((SELECT id FROM uploads WHERE upload_uuid = $1), $2, $3, $4, $5, $6, $7, $8)",
				uploadUUID,
				fileUUID,
				uploadUUID,
				filePart.FormName(),
				objKey,
				fileSize,
				mimeType,
				dateCreated,
			)
			if dbErr != nil {
				log.Printf("error from db exec call: %v", dbErr)
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, dbErr.Error(), http.StatusInternalServerError)
				encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": "could not register upload with user"})
				if encodeErr != nil {
					log.Printf("failed to write response: %v", encodeErr)
					return
				}
				return
			}
		} else {
			// Text field found
			_, _ = io.ReadAll(filePart)
		}
	}

	// Return upload UUID in response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(uploadUUID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
