package dto

type Agent struct {
	ID        int    `json:"id" db:"id"`
	AreaID    int    `json:"area_id" db:"area_id"`
	AgentName string `json:"agent_name" db:"agent_name"`
	AreaName  string `json:"area_name"`
}

type InsertAgent struct {
	ID        int    `db:"id"`
	AgentName string `json:"agent_name" db:"agent_name"`
	AreaName  string `json:"area_name" db:"area_name"`
}

type DashboardCount struct {
	Payment         *float64 `db:"SUM(payment)"`
	NumGallonsOrder *int     `db:"SUM(num_gallons_order)"`
	ReturnedGallons *int     `db:"SUM(returned_gallons)"`
}
