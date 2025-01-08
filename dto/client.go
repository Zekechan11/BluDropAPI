package dto

type InsertClient struct {
    FirstName string `json:"firstname" db:"firstname"`
    Lastname  string `json:"lastname" db:"lastname"`
    Email     string `json:"email" db:"email"`
    Username  string `json:"username" db:"username"`
    Password  string `json:"password" db:"password"`
    AreaId    string `json:"area_id" db:"area_id"`
    Role      string `json:"role" db:"role"`
}


type ClientModel struct {
	ClientId  int    `json:"client_id" db:"client_id"`
	FirstName string `json:"firstname" db:"firstname"`
	Lastname  string `json:"lastname" db:"lastname"`
	Email     string `json:"email" db:"email"`
	Username  string `json:"username" db:"username"`
	Password  string `json:"password" db:"password"`
	AreaId    string `json:"area_id" db:"area_id"`
	Role      string `json:"role" db:"role"`
}
