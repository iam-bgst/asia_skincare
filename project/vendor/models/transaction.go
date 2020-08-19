package models

import (
	"errors"
	"fmt"
	"forms"
	"time"

	"github.com/pborman/uuid"
)

type Transaction struct {
	Id       string               `json:"_id" bson:"_id,omitempty"`
	Product  []ProductTransaction `json:"product" bson:"product"`
	Paket    []Paket              `json:"paket" bson:"paket"`
	Discount []Discount           `json:"discount" bson:"discount"`
	Date     time.Time            `json:"date" bson:"date"`
	Delivery Delivery             `json:"delivery" bson:"delivery"`
	Subtotal int                  `json:"subtotal" bson:"subtotal"`
	Status   string               `json:"status" bson:"status"`
	To       To                   `json:"to" bson:"to"`
	From     From                 `json:"from" bson:"from"`
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

func (T *TransactionModel) Create(data forms.Transaction) (ret Transaction, err error) {
	id := uuid.New()

	// Getting the Product
	var prod []ProductTransaction
	discount_i := data.Product[0].Discount
	for i, p := range data.Product {
		if len(data.Product) == 0 {
			prod = nil
			fmt.Println("product kosong")
		} else {
			data_p, err1 := product_model.GetByMembership(p.Product, data.Membership)
			// fmt.Println(err1.Error())
			if err1 != nil {
				err = errors.New("Product id " + p.Product + " " + err1.Error())
				return
			}
			// fmt.Println(discount_i)
			if i > 0 {
				if p.Discount == discount_i {
					err = errors.New("discount used on product " + data.Product[i-1].Product)
					return
				}
			}
			data_d, err1 := discount_model.Get(p.Discount)
			if err1 != nil {
				err = errors.New("Discount Product id " + p.Discount + " " + err1.Error())
				return
			}
			prod = append(prod, ProductTransaction{
				Id:       data_p.Id,
				Name:     data_p.Name,
				Pricing:  data_p.Pricing.Price,
				Image:    data_p.Image,
				Discount: data_d,
			})
		}
	}

	// Getting the Paket
	var paket []Paket
	for _, p := range data.Paket {
		if len(data.Paket) == 0 {
			paket = nil
		} else {
			data_p, err1 := paket_model.Get(p)
			if err != nil {
				err = errors.New("Paket id " + p + " " + err1.Error())
				return
			}
			paket = append(paket, data_p)
		}

	}

	// Getting the Discount
	var dis []Discount
	for _, d := range data.Discount {
		if len(data.Discount) == 0 {
			dis = nil
		} else {
			data_discount, err1 := discount_model.Get(d)
			if err1 != nil {
				err = errors.New("Discount id " + d + " " + err1.Error())
				return
			}
			dis = append(dis, Discount{
				Id:           data_discount.Id,
				Discount:     data_discount.Discount,
				DiscountCode: data_discount.DiscountCode,
				Expired:      data_discount.Expired,
				Image:        data_discount.Image,
				Name:         data_discount.Name,
			})
		}
	}
	ret = Transaction{
		Id:       id,
		Date:     time.Now(),
		Discount: dis,
		Paket:    paket,
		Product:  prod,
	}

	// From
	ret.From.Address = data.From.Address
	ret.From.Name = data.From.Name
	ret.From.Number = data.From.Number

	// To
	ret.To.Address = data.To.Address
	ret.To.Name = data.To.Name
	ret.To.Number = data.To.Number

	// fmt.Println(ret)
	// err = db.Collection["transaction"].Insert(bson.M{
	// 	"_id":      id,
	// 	"product":  prod,
	// 	"paket":    paket,
	// 	"discount": dis,
	// 	"date":     time.Now(),
	// 	"from":     data.From,
	// 	"to":       data.To,
	// })
	return
}
