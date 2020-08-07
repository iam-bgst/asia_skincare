package forms

import "time"

type Transaction struct {
	Product  []Product `json:"product" bson:"product"`
	Paket    []Paket   `json:"paket" bson:"paket"`
	Date     time.Time `json:"date" bson:"date"`
	Delivery Delivery  `json:"delivery" bson:"delivery"`
	Subtotal int       `json:"subtotal" bson:"subtotal"`
	To       To        `json:"to" bson:"to"`
	From     From      `json:"from" bson:"from"`
}
type To struct {
	Name    string `json:"name" bson:"name"`
	Number  string `json:"number" bson:"number"`
	Address string `json:"address" bson:"address"`
}

type From struct {
	Name    string `json:"name" bson:"name"`
	Number  string `json:"number" bson:"number"`
	Address string `json:"address" bson:"address"`
}

type Delivery struct {
	Courier string `json:"courier" bson:"courier"`
	Resi    string `json:"resi" bson:"resi"`
	Price   string `json:"price" bson:"price"`
}
