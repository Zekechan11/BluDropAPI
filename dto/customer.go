package dto

type CustomerEntity struct {
	ID int `db:"id"`
	FirstName string `db:"firstname"`
	LastName string `db:"lastname"`
	Email string `db:"email"`
	Area string `db:"area"`
}