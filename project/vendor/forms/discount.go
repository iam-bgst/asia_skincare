package forms

import "time"

type Discount struct {
	Name         string    `json:"name"`
	Discount     int       `json:"discount"`
	DiscountCode string    `json:"discoundcode"`
	Image        string    `json:"image"`
	Expired      time.Time `json:"expired"`
}
