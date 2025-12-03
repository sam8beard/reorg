package models

type User struct {
	UserID Identity `json:"userID"`
	UserData Data   `json:"data"`
}

type Identity struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Data struct {
	Files       []FileData 	`json:"files"`
	Preferences Preferences `json:"preferences"`
}

type FileData struct { 
	Name string `json:"fileName"`
	Body []byte `json:"fileBody"`
}

type Preferences struct { 
} 
