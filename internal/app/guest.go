package app

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sam8beard/reorg/internal/models"
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
