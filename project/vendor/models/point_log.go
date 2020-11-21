package models

import "time"

type Point_log struct {
	Id      string    `json:"_id" bson:"_id,omitempty"`
	Date    time.Time `json:"date" bson:"date"`
	Desc    string    `json:"desc" bson:"desc"`
	Account Account2  `json:"account" bson:"account"`
}
