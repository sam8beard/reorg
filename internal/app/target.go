package app

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type TargetData struct {
	TargetUUID string `json:"targetUUID"`
	TargetName string `json:"targetName"`
}

func (s *Server) TargetHandler(w http.ResponseWriter, r *http.Request) {
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

	targetData := TargetData{}
	if err := json.Unmarshal(body, &targetData); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": "could not decode request body"})
		if encodeErr != nil {
			log.Printf("failed to write response: %v", encodeErr)
			return
		}
		return
	}
	// WORKING: properly receiving target info
	log.Printf("Target UUID: %s | Target Name: %s", targetData.TargetUUID, targetData.TargetName)

	// Insert row in targets table and retrieve newly created id
	var targetId int
	dbErr := s.DB.QueryRow(
		context.Background(),
		`
		INSERT INTO targets (target_uuid, name, user_id) 
		VALUES ($1, $2, NULL) 
		RETURNING id
		`,
		targetData.TargetUUID,
		targetData.TargetName,
	).Scan(&targetId)
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

	// Return target id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(targetId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
