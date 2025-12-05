package app

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type UploadForm struct {
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	uploadRoot := "./uploads"
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

	// use io.Copy ?
	// what do we want to return to the frontend?
	// the next logical step would be to flatten any directories and store/return a flattened list of all files

	// Read files in chunks
	for {
		part, err := multiReader.NextPart()
		if err == io.EOF {
			break
		}
		defer part.Close()
		if part.FileName() != "" {
			// File is found
			// Store files on disk temporarily
			filePath := filepath.Join(uploadRoot, part.FileName())
			os.MkdirAll(filepath.Dir(filePath), 0755)
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
}
