package app

import (
	"context"
	"encoding/json"
	"github.com/sam8beard/reorg/internal/models"
	"github.com/sam8beard/reorg/internal/rules"
	"io"
	"log"
	"net/http"
)

func (s *Server) PreviewHandler(w http.ResponseWriter, r *http.Request) {
	// Debugging
	log.Println("Preview endpoint firing")
	// Defer close request body
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

	ruleSet := models.RuleSet{}
	fileMetadata := make(map[string]models.FileMetadata, 0)

	// Unmarshal request body
	if err := json.Unmarshal(body, &ruleSet); err != nil {
		log.Printf("%v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": "could not decode request body"})
		if encodeErr != nil {
			log.Printf("failed to write response: %v", encodeErr)
			return
		}
		return
	}
	// Fetch metadata for each file from files table
	metadataRows, dbErr := s.DB.Query(
		context.Background(),
		"SELECT upload_id, id, name, size, mime_type, original_timestamp FROM files WHERE upload_id = $1",
		ruleSet.UploadID,
	)
	if dbErr != nil {
		log.Printf("error from db query call: %v", dbErr)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, dbErr.Error(), http.StatusInternalServerError)
		encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": "could not retrieve file metadata"})
		if encodeErr != nil {
			log.Printf("failed to write response: %v", encodeErr)
			return
		}
		return
	}

	// Scan all metadata rows and store
	for metadataRows.Next() {
		log.Println("Scanning row...")
		var md models.FileMetadata
		dbErr := metadataRows.Scan(
			&md.UploadID,
			&md.FileID,
			&md.FileName,
			&md.Size,
			&md.MimeType,
			&md.OGTimestamp,
		)
		if dbErr != nil {
			log.Printf("error scanning value from row: %v", dbErr)
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, dbErr.Error(), http.StatusInternalServerError)
			encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": "could not scan metadata row"})
			if encodeErr != nil {
				log.Printf("failed to write response: %v", encodeErr)
				return
			}
			return
		}

		// Debugging
		log.Printf("File metadata: %+v\n\n", md)

		// Store metadata for file
		fileID := md.FileID
		fileMetadata[fileID] = md
	}

	// Get evaluation result
	evalResult, err := rules.Evaluate(&ruleSet, fileMetadata)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		if encodeErr != nil {
			log.Printf("failed to write response: %v", encodeErr)
			return
		}
		return
	}
	//dummyFolders := make(map[string]*models.Folder, 0)
	//f := models.Folder{

	//	TargetUUID: "target-uuid",
	//	TargetName: "test-folder",
	//	Files: []models.File{
	//		{
	//			FileUUID: "file-uuid",
	//			FileName: "matched-file.pdf",
	//		},
	//	},
	//}

	//dummyFolders["target-uuid"] = &f

	//// Return dummy evaluation result for debugging
	//dummyResult := models.EvaluationResult{
	//	UploadUUID: "upload-uuid",
	//	Folders:    dummyFolders,
	//	Unmatched: models.UnmatchedFolder{
	//		Name: "unmatched",
	//		Files: []models.File{
	//			{
	//				FileUUID: "unmatched-file-uuid",
	//				FileName: "unmatched-file.png",
	//			},
	//		},
	//	},
	//}

	//// Return dummy result for debugging
	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	//if err := json.NewEncoder(w).Encode(dummyResult); err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

	// Debugging
	rules.LogEvalResult(evalResult)
	// Return evaluation result
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(evalResult); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
