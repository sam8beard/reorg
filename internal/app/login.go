package app

import (
	"encoding/json"
	"github.com/sam8beard/reorg/internal/models"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"net/http"
)

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {

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
		encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": "could not read request body"})
		if encodeErr != nil {
			log.Printf("failed to write response: %v", encodeErr)
			return
		}
		return
	}

	// Unmarshal request body
	loginRequest := models.LoginRequest{}
	if err := json.Unmarshal(body, &loginRequest); err != nil {
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

	// Fetch user record
	var hashedPassword string
	var userID string
	var username string
	var email string
	log.Println(loginRequest.UsernameOrEmail)
	dbErr := s.DB.QueryRow(
		r.Context(),
		`SELECT password_hash, id, username, email FROM users WHERE email=$1 OR username=$1`,
		loginRequest.UsernameOrEmail,
	).Scan(&hashedPassword, &userID, &username, &email)
	if dbErr != nil {
		log.Printf("error from db call: %v", dbErr)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, dbErr.Error(), http.StatusInternalServerError)
		encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": "could not verify user"})
		if encodeErr != nil {
			log.Printf("failed to write response: %v", encodeErr)
			return
		}
		return
	}
	// Validate password for user
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginRequest.Password)); err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, dbErr.Error(), http.StatusInternalServerError)
		encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": "username or password is incorrect"})
		if encodeErr != nil {
			log.Printf("failed to write response: %v", encodeErr)
			return
		}
		return
	}

	// Generate token
	token, err := s.JWTService.GenerateToken(string(userID))
	if err != nil {
		http.Error(w, "could not generate token", http.StatusInternalServerError)
		return
	}

	// Build response
	resp := models.AuthResponse{
		Token: token,
		User: models.User{
			ID:       userID,
			Username: username,
			Email:    email,
		},
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
