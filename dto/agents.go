package dto

type AgentsEntity struct {
	ID int `db:"id"`
	FirstName string `db:"firstname"`
	LastName string `db:"lastname"`
	Email string `db:"email"`
	Area string `db:"area"`
}

type AgentsModel struct {
	FirstName string `json:"firstname"`
	LastName string `json:"lastname"`
	Email string `json:"email"`
	Area string `json:"area"`
	Password string `json:"password"`
	Role string `json:"role"`
}