package dto

type InsertStaff struct {
	StaffId   int64    `json:"staff_id" db:"staff_id"`
	FirstName string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role      string `json:"role"`
	AreaId    int `json:"area_id" db:"area_id"`
}

type StaffModel struct {
	StaffId   int    `json:"staff_id" db:"staff_id"`
	FirstName string `json:"firstname" db:"firstname"`
	Lastname  string `json:"lastname" db:"lastname"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"password" db:"password"`
	Role      string `json:"role" db:"role"`
	AreaId    int 	 `json:"area_id" db:"area_id"`
	Area	  string `json:"area" db:"area"`
	Created	  string `json:"created" db:"created"`
}
