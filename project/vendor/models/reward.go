package models

type Rewards struct {
	Name       string `json:"name" bson:"name"`
	PricePoint int    `json:"pricepoint" bson:"pricepoint"`
}
