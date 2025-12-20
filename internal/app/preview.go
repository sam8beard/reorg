package app

import (
	"encoding/json"
	"github.com/sam8beard/reorg/internal/models"
	"github.com/sam8beard/reorg/internal/rules"
	"io"
	"log"
	"net/http"
)

func (s *Server) PreviewHandler(w http.ResponseWriter, r *http.Request) {
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

	// Unmarshal request body
	ruleSet := models.RuleSet{}
	if err := json.Unmarshal(body, &ruleSet); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": "could not decode request body"})
		if encodeErr != nil {
			log.Printf("failed to write response: %v", encodeErr)
			return
		}
		return
	}

	// Get evaluation result
	_, err = rules.Evaluate(&ruleSet)
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

	// Return dummy evaluation result for debugging
	dummyResult := models.EvaluationResult{
		UploadUUID: "upload-uuid",
		Folders: []models.Folder{
			{
				TargetUUID: "target-uuid",
				TargetName: "test-folder",
				Files: []models.File{
					{
						FileUUID: "file-uuid",
						FileName: "matched-file.pdf",
					},
				},
			},
		},
		Unmatched: models.UnmatchedFolder{
			Name: "unmatched",
			Files: []models.File{
				{
					FileUUID: "unmatched-file-uuid",
					FileName: "unmatched-file.png",
				},
			},
		},
	}

	// Return dummy result for debugging
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dummyResult); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//// Return evaluation result
	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	//if err := json.NewEncoder(w).Encode(evalResult); err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

}
