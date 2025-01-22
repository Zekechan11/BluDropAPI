package dto

type LoginData struct {
	Uid       int     `db:"uid"`
	FirstName string  `db:"firstname"`
	LastName  string  `db:"lastname"`
	UserName  *string `db:"username"`
	Area      *string `db:"area"`
	AreaId    *string `db:"area_id"`
	Type    *string `db:"type"`
	Email     string  `db:"email"`
	Role      string  `db:"role"`
	Password  string  `db:"password"`
}
