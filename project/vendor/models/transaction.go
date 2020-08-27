package models

import (
	"db"
	"errors"
	"fmt"
	"forms"
	"strconv"
	"strings"
	"time"

	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2/bson"
)

type Transaction struct {
	Id       string               `json:"_id" bson:"_id,omitempty"`
	Account  AccountTransaction   `json:"account" bson:"account"`
	Product  []ProductTransaction `json:"product" bson:"product"`
	Paket    []PaketTransaction   `json:"paket" bson:"paket"`
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

type TransactionModel struct{}

func (T *TransactionModel) Create(data forms.Transaction) (ret Transaction, err error) {
	id := uuid.New()

	// Get Account
	data_account, err1 := account_model.Get(data.Account)
	if err1 != nil {
		err = err1
		return
	}

	// Getting the Product
	var prod []ProductTransaction
	discount_i := data.Product[0].Discount
	for i, p := range data.Product {
		if len(data.Product) == 0 {
			prod = nil
		} else {
			_, err1 := account_model.GetDiscountUsed(data.Account, p.Discount)
			if err1 == nil {
				err = errors.New("discount is used on previous product")
				return
			}
			data_p, err1 := product_model.GetByMembership(p.Product, data.Membership)
			if err1 != nil {
				err = errors.New("Product id " + p.Product + " " + err1.Error())
				return
			}
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
				Qty:      p.Qty,
				Pricing:  data_p.Pricing.Price * p.Qty,
				Image:    data_p.Image,
				Discount: data_d,
			})
			err2 := account_model.AddDiscounUsed(data_account.Id, p.Discount)
			if err2 != nil {
				fmt.Println("log on line 96")
				err = err2
				return
			}
		}
	}

	// Getting the Paket
	var paket []PaketTransaction
	for _, p := range data.Paket {
		if len(data.Paket) == 0 {
			paket = nil
		} else {
			data_p, err1 := paket_model.GetByMembership(p.Paket, data.Membership)
			if err != nil {
				err = errors.New("Paket id " + p.Paket + " " + err1.Error())
				return
			}
			paket = append(paket, PaketTransaction{
				Id:      data_p.Id,
				Image:   data_p.Image,
				Qty:     p.Qty,
				Name:    data_p.Name,
				Point:   data_p.Point,
				Pricing: data_p.Pricing.Price * p.Qty,
				Product: data_p.Product,
				Stock:   data_p.Stock,
			})
		}
	}

	// Getting the Discount
	var dis []Discount
	for _, d := range data.Discount {
		_, err1 := account_model.GetDiscountUsed(data.Account, d)
		if err1 == nil {
			err = errors.New("discount is used on your account")
			return
		}
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
			err2 := account_model.AddDiscounUsed(data.Account, d)
			if err2 != nil {
				fmt.Println("log on line 145")
				err = err2
				return
			}
		}
	}
	ret = Transaction{
		Id:       id,
		Date:     time.Now(),
		Discount: dis,
		Paket:    paket,
		Product:  prod,
		Account:  data_account,
	}

	// From
	ret.From.Address = data.From.Address
	ret.From.Name = data.From.Name
	ret.From.Number = data.From.Number

	// To
	ret.To.Address = data.To.Address
	ret.To.Name = data.To.Name
	ret.To.Number = data.To.Number

	// Proses Subtotal
	subtotal := 0
	total_discount := 0
	total := 0
	for _, dis := range prod {
		total_discount += dis.Discount.Discount
		subtotal += dis.Pricing
	}
	for _, pak := range paket {
		subtotal += pak.Pricing
	}
	for _, dis := range dis {
		total_discount += dis.Discount
	}

	// log.Println(subtotal)
	discount := (subtotal * total_discount) / 100
	// log.Println(discount)
	total = subtotal - discount
	// log.Println(total)
	ret.Subtotal = total

	// Insert into mongo
	err = db.Collection["transaction"].Insert(bson.M{
		"_id":      id,
		"date":     time.Now(),
		"account":  data_account,
		"product":  prod,
		"paket":    paket,
		"discount": dis,
		"from":     data.From,
		"to":       data.To,
		"delivery": bson.M{
			"courier": data.Delivery.Courier,
			"service": data.Delivery.Service,
			"resi":    "",
			"price":   data.Delivery.Price,
			"code":    data.Delivery.Code,
		},
	})
	if err != nil {
		return
	}

	return
}

func (T *TransactionModel) UpdateStatus(id, status string) (err error) {
	err = db.Collection["transaction"].Update(bson.M{
		"_id": id,
	}, bson.M{
		"status": status,
	})
	return
}

func (T *TransactionModel) UpdateResi(id, resi string) (err error) {
	err = db.Collection["transaction"].Update(bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			"delivery.resi": resi,
		},
	})
	return
}

func (T *TransactionModel) HistoyTransaction(id_account, filter, sort, pageNo, perPage string) (data []Transaction, count int, err error) {
	sorting := sort
	if strings.Contains(sort, "asc") {
		sorting = strings.Replace(sort, "|asc", "", -1)
	} else if strings.Contains(sort, "desc") {
		sorting = strings.Replace(sort, "|desc", "", -1)
		sorting = "-" + sorting
	}
	regex := bson.M{"$regex": bson.RegEx{Pattern: filter, Options: "i"}}
	pn, _ := strconv.Atoi(pageNo)
	pp, _ := strconv.Atoi(perPage)
	err = db.Collection["transaction"].Find(bson.M{
		"account._id": id_account,
		"$or": []interface{}{
			bson.M{"name": regex},
		},
	}).Sort(sorting).Skip((pn - 1) * pp).Limit(pp).All(&data)
	if err != nil {
		return
	}
	count, err = db.Collection["transaction"].Find(bson.M{"account._id": id_account}).Count()
	if err != nil {
		return
	}
	return
}
