package dto

type ClientOrder struct {
	TotalPrice float64 `db:"total_price"`
	NumGallons int     `db:"num_gallons_order"`
	AreaID     int     `db:"area_id"`
}

type COL struct {
	ExistingRecord     int  `db:"COUNT(*)"`
	PreviousNumGallons *int `db:"total_containers_on_loan"`
}
