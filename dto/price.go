package dto

type PricingModel struct {
	PricingId int     `db:"pricing_id" json:"pricing_id"`
	Dealer    float64 `db:"dealer" json:"dealer"`
	Regular   float64  `db:"regular" json:"regular"`
}

type InsertPricing struct {
	Dealer  float64 `db:"dealer" json:"dealer"`
	Regular float64  `db:"regular" json:"regular"`
}
