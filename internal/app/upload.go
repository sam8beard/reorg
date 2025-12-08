package app

import (
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type UploadForm struct {
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Generate new upload ID
	uploadId := uuid.New()
	uploadRoot := "./uploads/" + uploadId.String()

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
		if err == io.EOF {
			break
		}
		defer part.Close()
		if part.FileName() != "" {
			// File is found
			// Store files on disk
			filePath := filepath.Join(uploadRoot, part.FileName())
			os.MkdirAll(filepath.Dir(filePath), 0755)

			// this is the full path to the file
			// need to find a way to reconstruct directory path on the backend...?
			log.Printf("file path: %s", filePath)
			// Create file on server to store file body
			out, _ := os.Create(filePath)
			io.Copy(out, part)
			out.Close()
		} else {
			// Text field found
			data, _ := io.ReadAll(part)
			log.Printf("field: %s \t value: %s\n", part.FormName(), string(data))
		}
	}

	// Return upload id in response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(uploadId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
