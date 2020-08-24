package models

type Rewards struct {
	Name       string `json:"name" bson:"name"`
	PricePoint int    `json:"pricepoint" bson:"pricepoint"`
	Reward     string `json:"reward" bson:"reward"`
	Image      string `json:"image" bson:"image"`
}
