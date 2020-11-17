package forms

import "time"

type Transaction struct {
	Product  []ProductTransaction `json:"product" bson:"product"`
	Discount []string             `json:"discount" bson:"discount"`
	Delivery Delivery             `json:"delivery" bson:"delivery"`
	Payment  string               `json:"payment"`
	Subtotal int                  `json:"subtotal" bson:"subtotal"`
	To       To                   `json:"to" bson:"to"`
	From     From                 `json:"from" bson:"from"`
}
type Evidence struct {
	Total   int       `json:"total"`
	Name    string    `json:"name"`
	Send_by string    `json:"send_by"`
	Time    time.Time `json:"time"`
	Image   string    `json:"image"`
}

type PaketTransaction struct {
	Paket string `json:"paket"`
	Qty   int    `json:"qty"`
}

type ProductTransaction struct {
	Product  string `json:"product"`
	Qty      int    `json:"qty"`
	Discount string `json:"discount"`
}

type To struct {
	Account string `json:"account"`
	Address string `json:"address" bson:"address"`
}

type From struct {
	Account string `json:"account"`
}

type Delivery struct {
	Courier string `json:"courier"`
	Service string `json:"service"`
	Resi    string `json:"resi"`
	Price   string `json:"price"`
	Code    string `json:"code"`
}
