package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/sam8beard/reorg/internal/auth/middleware"
	"github.com/schollz/progressbar/v3"
	"io"
	"log"
	"mime"
	"mime/multipart"
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
	// Determine whether request came from a guest or registered user
	userID := r.Context().Value(middleware.CtxKeyUserID).(string)
	isGuest := r.Context().Value(middleware.CtxKeyGuest).(bool)

	// Should print true or guest
	log.Println(r.Context().Value(middleware.CtxKeyGuest))
	log.Println(r.Context())

	// Generate new upload UUID
	uploadID := uuid.New()

	// Close request body
	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			log.Fatalf("could not close request body: %v", closeErr)
		}
	}()

	// Get upload size for progress bar
	uploadSize := r.ContentLength

	// Build progress bar
	pBar := progressbar.DefaultBytes(uploadSize, "Uploading files...")

	// Wrap request body
	pReader := io.TeeReader(r.Body, pBar)

	// Get boundary
	_, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, err.Error(), http.StatusBadRequest)
		encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": "could not parse request body"})
		if encodeErr != nil {
			log.Printf("failed to write response: %v", encodeErr)
			return
		}
		return
	}

	// Create multipart reader to stream files and update progress bar
	multiReader := multipart.NewReader(pReader, params["boundary"])

	// Build query on whether upload was made my registered user or guest
	var uploadQuery string
	if isGuest {
		uploadQuery = "INSERT INTO uploads (id, guest_id) VALUES ($1, $2)"
	} else {

		uploadQuery = "INSERT INTO uploads (id, user_id) VALUES ($1, $2)"
	}

	// Insert row in uploads table
	_, dbErr := s.DB.Exec(
		context.Background(),
		uploadQuery,
		uploadID,
		userID,
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

		// filePart.FormName() = fileName, filePart.FileName() = dateCreated
		if filePart.FormName() != "" && filePart.FormName() != ".DS_Store" {
			// File is found

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
			}
			dateMilliInt64 := int64(dateMilliInt)
			dateCreated := time.UnixMilli(dateMilliInt64).UTC()
			// Generate file uuid
			fileID := uuid.New()
			// Get mime type
			mimeType := filePart.Header.Get("Content-Type")
			// Build user metadata
			userMetadata := map[string]string{
				"original-file-name": filePart.FormName(),
				"last-modified":      dateCreated.String(),
				"mime-type":          string(mimeType),
				"upload-timestamp":   time.Now().String(),
				"file-id":            fileID.String(),
				"upload-id":          uploadID.String(),
				//"user-id": userID // Once auth is set up
			}
			// Get content disposition
			contentDisposition := fmt.Sprintf("attachment; filename=\"%s\"", filePart.FormName())

			// Make upload progress bar to display in logs

			/*
				Should we delete all files after a certain date? I guess there is no point in promising persistence.


			*/
			opts := minio.PutObjectOptions{
				UserMetadata:       userMetadata,
				ContentType:        string(mimeType),
				ContentDisposition: contentDisposition,
			}
			objKey := fmt.Sprintf("%s/%s_%s", uploadID, fileID, filePart.FormName())
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

			// Get file size
			fileSize := uploadInfo.Size

			// Insert row in files table
			_, dbErr := s.DB.Exec(
				context.Background(),
				"INSERT INTO files (upload_id, id, name, s3_key, size, mime_type, original_timestamp) VALUES ((SELECT id FROM uploads WHERE id = $1), $2, $3, $4, $5, $6, $7)",
				uploadID,
				fileID,
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
			filePart.Close()
		} else {
			// Text field found
			_, _ = io.ReadAll(filePart)
		}
	}

	// Finish progress bar and display success message
	if err := pBar.Finish(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, dbErr.Error(), http.StatusInternalServerError)
		encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": "upload incomplete due to network error or file corruption"})
		if encodeErr != nil {
			log.Printf("failed to write response: %v", encodeErr)
			return
		}
		log.Printf("Upload Unsuccessful for upload ID:%s, %v", uploadID.String(), err)
		return
	}
	log.Printf("Upload Successful for upload ID:%s", uploadID.String())

	// Return upload ID in response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(uploadID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
