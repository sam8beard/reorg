package app

import (
	"encoding/json"
	"github.com/google/uuid"
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
	UploadID     string        `json:"uploadID"`
	Files        []models.File `json:"files"`
	Targets      []target      `json:"targets"`
	RuleBindings []ruleBinding `json:"ruleBindings"`
}

//type ruleSet models.EvaluationResult

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

	// Data for entries
	uploadID := orgData.RuleSet.UploadID
	rules := orgData.Rules
	ruleBindings := orgData.RuleBindings
	ruleSet := orgData.RuleSet

	// Debugging
	log.Printf("\n\nRules in download: \n%+v", rules)
	log.Printf("\n\nRule bindings in download: \n%+v", ruleBindings)
	log.Printf("\n\nRule set in download: \n%+v", ruleSet)

	// Print all rules
	for _, rule := range rules {
		log.Printf("%+v", rule.Conditions)
		log.Printf("%s", rule.RuleName)
		log.Printf("%s", rule.RuleID)
	}

	// Begin transaction
	tx, err := s.DB.Begin(r.Context())
	if err != nil {
		log.Printf("%v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		if encodeErr != nil {
			log.Printf("failed to write response: %v", encodeErr)
			return
		}
		return
	}
	// Ensure rollback on failure
	defer func() {
		if err != nil {
			_ = tx.Rollback(r.Context())
		}
	}()

	// Write ruleset and get newly created ruleset ID
	var ruleSetID string
	err = tx.QueryRow(
		r.Context(),
		`INSERT INTO rulesets 
		(id, user_id, ruleset_json, upload_id)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		uuid.New(),
		userID,
		ruleSet,
		uploadID,
	).Scan(&ruleSetID)
	if err != nil {
		log.Printf("%v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		if encodeErr != nil {
			log.Printf("failed to write response: %v", encodeErr)
			return
		}
		return
	}

	// Write all rules
	for _, rule := range rules {
		_, err := tx.Exec(
			r.Context(),
			`INSERT INTO rules
			(id, user_id, name, conditions_json)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (user_id, name, conditions_json) DO NOTHING`,
			rule.RuleID,
			userID,
			rule.RuleName,
			rule.Conditions,
		)
		if err != nil {
			log.Printf("%v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			if encodeErr != nil {
				log.Printf("failed to write response: %v", encodeErr)
				return
			}
			return
		}
	}

	// Build target map for id -> name reference
	targetMap := make(map[string]string, 0)
	for _, target := range ruleSet.Targets {
		targetMap[target.TargetID] = target.TargetName
	}

	// Write all rule bindings
	for _, binding := range ruleBindings {
		_, err := tx.Exec(
			r.Context(),
			`INSERT INTO rule_bindings
			(id, ruleset_id, rule_id, target_id, target_name)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (ruleset_id, target_id, rule_id) DO NOTHING`,
			uuid.New(),
			ruleSetID,
			binding.RuleID,
			binding.TargetID,
			targetMap[binding.TargetID],
		)
		if err != nil {
			log.Printf("%v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			if encodeErr != nil {
				log.Printf("failed to write response: %v", encodeErr)
				return
			}
			return
		}
	}

	// Commit transaction
	err = tx.Commit(r.Context())
	if err != nil {
		log.Printf("%v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		if encodeErr != nil {
			log.Printf("failed to write response: %v", encodeErr)
			return
		}
		return
	}

}
