package dto

type PricingModel struct {
	PricingId int `db:"pricing_id" json:"pricing_id"`
	Value float64 `db:"value" json:"value"`
	Type string `db:"type" json:"type"`
}

type InsertPricing struct {
	Value float64 `db:"value" json:"value"`
	Type string `db:"type" json:"type"`
}