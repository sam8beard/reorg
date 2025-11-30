package models

type User struct {
	UserID Identity `json:"userID"`
}

type Identity struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
