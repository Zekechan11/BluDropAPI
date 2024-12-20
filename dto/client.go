package dto

type InsertClient struct {
	FirstName string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Username     string `json:"username"`
	Password  string `json:"password"`
	Address string `json:"address"`
	Role      string `json:"role"`
}

type ClientModel struct {
	ClientId   int    `json:"client_id" db:"client_id"`
	FirstName string `json:"firstname" db:"firstname"`
	Lastname  string `json:"lastname" db:"lastname"`
	Email     string `json:"email" db:"email"`
	Username     string `json:"username" db:"username"`
	Password  string `json:"password" db:"password"`
	Address string `json:"address" db:"address"`
	Role      string `json:"role" db:"role"`
}