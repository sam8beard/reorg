package app

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sam8beard/reorg/internal/models"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"net/http"
	"regexp"
)

func (s *Server) SignupHandler(w http.ResponseWriter, r *http.Request) {
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
	signupRequest := models.SignupRequest{}
	if err := json.Unmarshal(body, &signupRequest); err != nil {
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

	/* Validate signup request */

	// Email validation
	emailR, _ := regexp.Compile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	validEmail := emailR.MatchString(signupRequest.Email)

	// Username validation
	var validUsername bool
	var validUsernameLen bool
	var validUsernameReg bool
	// Length is between 4 and 15 (inclusive)
	if len(signupRequest.Username) >= 4 || len(signupRequest.Username) <= 15 {
		validUsernameLen = true
	}
	// Contains only letters and numbers
	userRFirst, _ := regexp.Compile(`^[A-Za-z0-9]+$`)
	// Does not contain only numbers
	userRSecond, _ := regexp.Compile(`^\d+$`)
	validUsernameReg = userRFirst.MatchString(signupRequest.Username) && !userRSecond.MatchString(signupRequest.Username)

	// Password validation
	var validPass bool
	var validPassLen bool
	var validPassReg bool
	// Length is at least 8
	if len(signupRequest.Password) >= 8 {
		validPassLen = true
	}
	// Contains at least one digit
	passRFirst, _ := regexp.Compile(`\d`)
	// Contains at least one symbol
	passRSecond, _ := regexp.Compile(`\W`)
	validPassReg = passRFirst.MatchString(signupRequest.Password) && passRSecond.MatchString(signupRequest.Password)

	// Results of username and password validation
	validUsername = validUsernameLen && validUsernameReg
	validPass = validPassLen && validPassReg

	// Return error indicating invalid signup information
	if !validEmail || !validUsername || !validPass {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": "invalid signup input"})
		if encodeErr != nil {
			log.Printf("failed to write response: %v", encodeErr)
			return
		}
		return
	}

	// Verify new user is not already registered
	var exists bool
	dbErr := s.DB.QueryRow(
		r.Context(),
		`SELECT EXISTS (
			SELECT 1 
			FROM users 
			WHERE email = $1 OR username = $2
		)`,
		signupRequest.Email,
		signupRequest.Username,
	).Scan(&exists)
	if dbErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": dbErr.Error()})
		if encodeErr != nil {
			log.Printf("failed to write response: %v", encodeErr)
			return
		}
		return
	}
	// Account already exists
	if exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": "account already exists"})
		if encodeErr != nil {
			log.Printf("failed to write response: %v", encodeErr)
			return
		}
		return
	}

	// Generate hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("%v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": "could not generate encoded user information"})
		if encodeErr != nil {
			log.Printf("failed to write response: %v", encodeErr)
			return
		}
		return
	}

	// Create user in database and get new user ID
	userID := uuid.New().String()
	_, dbErr = s.DB.Exec(
		context.Background(),
		`INSERT INTO users (id, username, email, password_hash) VALUES ($1, $2, $3, $4)`,
		userID,
		signupRequest.Username,
		signupRequest.Email,
		hashedPassword,
	)

	// TODO: HANDLE ACCOUNT ALREADY EXISTS ERROR
	if dbErr != nil {
		log.Printf("error from db exec call: %v", dbErr)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, dbErr.Error(), http.StatusInternalServerError)
		encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": "could not register user"})
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
			Username: signupRequest.Username,
			Email:    signupRequest.Email,
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
