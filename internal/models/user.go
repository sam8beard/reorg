package models

import ()

type User struct {
	ID       int    `json:"userID"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
type Identity struct {
	Type       string // "user" | "guest"
	UserID     string // user ID or guest ID
	SesssionID string // JWT
}

type Data struct {
	Files       []FileData  `json:"files"`
	Preferences Preferences `json:"preferences"`
}

type FileData struct {
	Name string `json:"fileName"`
	Body []byte `json:"fileBody"`
}

type Preferences struct {
}
