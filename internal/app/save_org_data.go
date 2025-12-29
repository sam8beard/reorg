package app

import (
	"encoding/json"
	"github.com/sam8beard/reorg/internal/auth/middleware"
	"github.com/sam8beard/reorg/internal/models"
	"io"
	"log"
	"net/http"
)

type ruleBinding struct {
	RuleID   string `json:"ruleID"`
	TargetID string `json:"targetID"`
}

type target struct {
	TargetID   string `json:"targetID"`
	TargetName string `json:"targetName"`
}

type ruleSet struct {
	UploadID string        `json:"uploadID"`
	Files    []models.File `json:"files"`
	Targets  []target      `json:"targets"`
}

type OrgData struct {
	Rules        []models.Rule `json:"rules"`
	RuleBindings []ruleBinding `json:"ruleBindings"`
	RuleSet      ruleSet       `json:"ruleSet"`
}

func (s *Server) SaveOrg(w http.ResponseWriter, r *http.Request) {

	// Defer close request body
	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			log.Fatalf("could not close request body: %v", closeErr)
		}
	}()

	// Get user ID
	userID := r.Context().Value(middleware.CtxKeyUserID).(string)
	log.Printf("User ID: %s", userID)

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
	orgData := OrgData{}
	if err := json.Unmarshal(body, &orgData); err != nil {
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

	// Debugging
	log.Printf("%+v", orgData)

	// Data for entries
	//rules := orgData.Rules
	//ruleBindings := orgData.RuleBindings
	//ruleSet := orgData.RuleSet

	// Write all entries to database
	log.Println("WRITE ENTRIES HERE")

}
