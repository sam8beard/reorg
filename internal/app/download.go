package app

import (
	"archive/zip"
	"context"
	"encoding/json"
	"github.com/minio/minio-go/v7"
	"github.com/sam8beard/reorg/internal/models"
	"github.com/sam8beard/reorg/internal/rules"
	"io"
	"log"
	"net/http"
	"path/filepath"
)

func (s *Server) DownloadZipHandler(w http.ResponseWriter, r *http.Request) {
	// Defer close on request body
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
		err := json.NewEncoder(w).Encode(map[string]string{"error": "could not read request body"})
		if err != nil {
			log.Printf("failed to write response: %v", err)
			return
		}
		return
	}

	fileStructure := models.EvaluationResult{}

	// Unmarshal request body
	if err := json.Unmarshal(body, &fileStructure); err != nil {
		log.Printf("%v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(map[string]string{"error": "could not decode request body"})
		if err != nil {
			log.Printf("failed to write response: %v", err)
			return
		}
		return
	}

	// Debugging
	rules.LogEvalResult(&fileStructure)

	// Get upload UUID
	uploadUUID := fileStructure.UploadUUID

	// Keep track of folder names to folder UUIDs for db retrieval
	folderMap := make(map[string]string, 0)
	fileMap := make(map[string]string, 0)

	// Add all folders to zip archive
	folders := fileStructure.Folders
	for _, folder := range folders {

		targetUUID := folder.TargetUUID
		folderName := folder.TargetName

		// Track folder
		folderMap[targetUUID] = folderName

		// Map file UUID to target UUID
		for _, file := range folder.Files {
			fileMap[file.FileUUID] = targetUUID
		}
	}

	// Get file metadata needed for download
	opts := minio.ListObjectsOptions{
		Prefix:       uploadUUID,
		WithMetadata: true,
		Recursive:    true,
	}

	// Set headers for downloadable zip archive
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="organized_files.zip"`)

	// Zip writer
	zipWriter := zip.NewWriter(w)

	// List all objects for this upload
	for obj := range s.Minio.ListObjects(
		context.Background(),
		s.MinioBucket,
		opts,
	) {

		log.Printf("%+v", obj)

		// Object key
		key := obj.Key

		// Debugging
		log.Println(key)

		// File name
		fileName := obj.UserMetadata["X-Amz-Meta-Original-File-Name"]

		// File UUID
		fileUUID := obj.UserMetadata["X-Amz-Meta-File-Uuid"]

		// Get target UUID that we mapped to this file
		targetUUID := fileMap[fileUUID]
		// Get folder name that we mapped to this target UUID
		folderName := folderMap[targetUUID]

		// Download file body
		getOpts := minio.GetObjectOptions{}
		fileBody, err := s.Minio.GetObject(
			context.Background(),
			s.MinioBucket,
			key,
			getOpts,
		)
		if err != nil {
			log.Printf("%v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(map[string]string{"error": "could not download files from object storage"})
			if err != nil {
				log.Printf("failed to write response: %v", err)
				return
			}
			return
		}

		// Create entry in archive
		filePath := filepath.Join(folderName, fileName)

		// Debuggin
		log.Printf("Folder name: %s", folderName)
		log.Printf("File name: %s", fileName)
		log.Printf("File path: %s", filePath)

		// Create file in zip archive
		entry, err := zipWriter.Create(filePath)
		if err != nil {
			log.Fatalf("%v", err)
		}

		info, err := fileBody.Stat()
		if err != nil {
			log.Printf("stat error: %v", err)
			continue
		}

		// Debugging
		log.Printf("%v", info.Size)
		log.Printf("%v", info.Metadata)

		// Stream file body to archive
		_, err = io.Copy(entry, fileBody)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// Close file body
		if err := fileBody.Close(); err != nil {
			log.Fatalf("%v", err)
		}

	}

	/*
		DONT FORGET TO WRITE UNSORTED FILES TOO
	*/

	err = zipWriter.Close()
	if err != nil {
		log.Fatalf("%v", err)
	}

}
