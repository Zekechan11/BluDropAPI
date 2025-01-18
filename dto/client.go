package dto

type InsertClient struct {
	ClientId  int64  `json:"client_id" db:"client_id"`
	FirstName string `json:"firstname" db:"firstname"`
	Lastname  string `json:"lastname" db:"lastname"`
	Email     string `json:"email" db:"email"`
	Username  string `json:"username" db:"username"`
	Password  string `json:"password" db:"password"`
	AreaId    int    `json:"area_id" db:"area_id"`
	Role      string `json:"role" db:"role"`
}

type ClientModel struct {
	ClientId      int      `json:"client_id" db:"client_id"`
	FirstName     string   `json:"firstname" db:"firstname"`
	Lastname      string   `json:"lastname" db:"lastname"`
	Email         string   `json:"email" db:"email"`
	Username      string   `json:"username" db:"username"`
	Password      string   `json:"password" db:"password"`
	AreaId        string   `json:"area_id" db:"area_id"`
	Area          string   `json:"area" db:"area"`
	Role          string   `json:"role" db:"role"`
	Status        string   `json:"status" db:"status"`
	Type          string   `json:"type" db:"type"`
	ContainerLoan *int     `json:"total_containers_on_loan" db:"total_containers_on_loan"`
	TotalPayable  *float64 `json:"total_payable" db:"total_payable"`
	Created       string   `json:"created" db:"created"`
}
