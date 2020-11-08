package forms

import "time"

type Discount struct {
	Name         string `json:"name"`
	Discount     int    `json:"discount"`
	DiscountCode string `json:"discountcode"`
	// Image        string    `json:"image"`
	StartAt time.Time `json:"startAt"`
	EndAt   time.Time `json:"endAt"`
}
