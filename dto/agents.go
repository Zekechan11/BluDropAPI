package dto

type AgentsEntity struct {
    ID        int    `db:"id" json:"ID"`
    FirstName string `db:"firstname" json:"FirstName"`
    LastName  string `db:"lastname" json:"LastName"`
    Email     string `db:"email" json:"Email"`
    Area      string `db:"area" json:"Area"`
}

type AgentsModel struct {
	FirstName string `json:"firstname"`
	LastName string `json:"lastname"`
	Email string `json:"email"`
	Area string `json:"area"`
	Password string `json:"password"`
	Role string `json:"role"`
}