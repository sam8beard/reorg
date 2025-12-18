package app

import (
	"encoding/json"
	"github.com/sam8beard/reorg/internal/models"
	"io"
	"log"
	"net/http"
)

func (s *Server) RuleHandler(w http.ResponseWriter, r *http.Request) {
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

	ruleJson := models.Rule{}
	if err := json.Unmarshal(body, &ruleJson); err != nil {
		log.Printf("unmarshal err: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": "could not decode request body"})
		if encodeErr != nil {
			log.Printf("failed to write response: %v", encodeErr)
			return
		}
		return
	}

	// Pretty printing rule struct for testing
	jsonRule, _ := json.MarshalIndent(ruleJson, "", " ")
	log.Println(string(jsonRule))

	// Gotta decide whether or not a separate table for rules vs rule sets
	// Could use ruleset id as a foreign key in each rule table row
}
