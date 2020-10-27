package models

import "time"

type Reedem struct {
	Id     string    `json:"_id" bson:"_id"`
	Code   string    `json:"code" bson:"code"`
	Reward Rewards   `json:"reward" bson:"reward"`
	Date   time.Time `json:"date" bson:"date"`
	Valid  bool      `json:"valid" bson:"valid"`
}
