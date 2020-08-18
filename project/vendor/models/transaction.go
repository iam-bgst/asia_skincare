package models

import (
	"db"
	"forms"
	"time"

	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2/bson"
)

type Transaction struct {
	Id       string     `json:"_id" bson:"_id,omitempty"`
	Product  []Product  `json:"product" bson:"product"`
	Paket    []Paket    `json:"paket" bson:"paket"`
	Discount []Discount `json:"discount" bson:"discount"`
	Date     time.Time  `json:"date" bson:"date"`
	Delivery Delivery   `json:"delivery" bson:"delivery"`
	Subtotal int        `json:"subtotal" bson:"subtotal"`
	Status   string     `json:"status" bson:"status"`
	To       To         `json:"to" bson:"to"`
	From     From       `json:"from" bson:"from"`
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

type TransactionModel struct{}

func (T *Transaction) Create(data forms.Transaction) (err error) {
	id := uuid.New()

	// Getting the Product
	var product []Product1
	for _, p := range data.Product {
		data_p, err := product_model.Get(p)
		if err != nil {
			break
		}
		product = append(product, data_p)
	}

	// Getting the Paket
	var paket []Paket
	for _, p := range data.Paket {
		data_p, err := paket_model.Get(p)
		if err != nil {
			break

		}
		paket = append(paket, data_p)
	}

	// Getting the Discount
	var discount []Discount
	for _, d := range data.Discount {
		data_discount, err := discount_model.Get(d)
		if err != nil {
			break
		}
		discount = append(discount, data_discount)
	}
	err = db.Collection["transaction"].Insert(bson.M{
		"_id":      id,
		"product":  product,
		"paket":    paket,
		"discount": discount,
		"date":     time.Now(),
	})
	return
}
