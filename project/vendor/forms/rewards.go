package forms

import "time"

type Rewards struct {
	Name       string    `json:"name"`
	PricePoint int       `json:"pricepoint"`
	Reward     string    `json:"reward"`
	Desc       string    `json:"desc"`
	Image      string    `json:"image"`
	Start      time.Time `json:"start"`
	End        time.Time `json:"end"`
}
