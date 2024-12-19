package dto

type CustomerEntity struct {
	ID        int    `db:"id" json:"ID"`
	FirstName string `db:"firstname" json:"FirstName"`
	LastName  string `db:"lastname" json:"LastName"`
	Email     string `db:"email" json:"Email"`
	Area      string `db:"area" json:"Area"`
}
