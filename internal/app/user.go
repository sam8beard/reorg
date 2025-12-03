package app

import (
	"encoding/json"
	"fmt"
	"github.com/sam8beard/reorg/internal/models"
	"io"
	"log"
	"net/http"
)

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//type UserData struct {
//}

// dummy data for testing
var dummyUsername = "username"
var dummyPassword = "password"

var dummyUsername2 = "muel"
var dummyPassword2 = "angel"

/*
Get user data on login
*/
func UserHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			log.Fatalf("could not close request body: %v", closeErr)
		}
	}()

	/*
		Here is where we would typically query the database our user data is stored in.

		For now, we'll just return a dummy user.
	*/

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
	loginForm := LoginForm{}
	if err := json.Unmarshal(body, &loginForm); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": "could not decode request body"})
		if encodeErr != nil {
			log.Printf("failed to write response: %v", encodeErr)
			return
		}
		return
	}

	// Validate user (dummy validation for now)
	username := loginForm.Username
	password := loginForm.Password

	// If user exists, fetch user data from database
	// userData, err := user.FetchUserData()

	// Dummy validation
	userData, err := func(u string, p string) (models.User, error) {
		// Eventually, we will send a request to our database
		// to validate the user
		if u == dummyUsername && p == dummyPassword {
			validUser := models.User{
				UserID: models.Identity{
					ID:       1234,
					Username: username,
					Email:    "user@user.com",
				},
			}
			return validUser, nil
		// Second dummy user
		} else if u == dummyUsername2 && p == dummyPassword2 { 
			
			validUser := models.User{
				UserID: models.Identity{
					ID:       5678,
					Username: username,
					Email:    "user@user.com",
				},
			}
			return validUser, nil

		} 
		return models.User{}, fmt.Errorf("user does not exist")
	}(username, password)

	// User does not exist
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := map[string]string{
			"error": "user does not exist",
		}
		encodeErr := json.NewEncoder(w).Encode(errMsg)
		if encodeErr != nil {
			log.Printf("failed to write response: %v", encodeErr)
			return
		}
		return
	}

	// For now, return dummy data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(userData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
