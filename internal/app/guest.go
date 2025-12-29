package app

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sam8beard/reorg/internal/models"
	"log"
	"net/http"
)

func (s *Server) GuestHandler(w http.ResponseWriter, r *http.Request) {

	// Generate guest ID and guest token
	guestID := uuid.New().String()
	token, _ := s.JWTService.GenerateGuestToken(guestID)

	// Build response
	resp := models.AuthGuestResponse{
		Token:   token,
		GuestID: guestID,
	}

	log.Println(token)
	log.Println(guestID)
	log.Println(resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
