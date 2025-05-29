package dto

type Profile struct {
	ID        int    `json:"id" db:"id"`
	Firstname string `json:"firstname" db:"firstname"`
	Lastname  string `json:"lastname" db:"lastname"`
	Email     string `json:"email" db:"email"`
}

type Password struct {
	ID             int    `json:"id" db:"id"`
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"password" db:"password"`
}