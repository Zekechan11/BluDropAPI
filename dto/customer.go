package dto

type CustomerEntity struct {
	ID        int    `db:"id" json:"ID"`
	FirstName string `db:"firstname" json:"FirstName"`
	LastName  string `db:"lastname" json:"LastName"`
	Email     string `db:"email" json:"Email"`
	Area      string `db:"area" json:"Area"`
}

type CustomerTransaction struct {
    ID               int     `db:"id" json:"id"`
    NumGallonsOrder  int     `db:"num_gallons_order" json:"num_gallons_order"`
    ReturnedGallons  int     `db:"returned_gallons" json:"returned_gallons"`
    Date             string  `db:"date" json:"date"`
    Payment          float64 `db:"payment" json:"payment"`
}
