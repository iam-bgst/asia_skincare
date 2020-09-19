package models

import (
	"addon"
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

const (
	NOTPAYED = iota
	PACKED
	SENT
	DONE
	CENCELED
)

type Transaction struct {
	Id          string               `json:"_id" bson:"_id,omitempty"`
	Product     []ProductTransaction `json:"product" bson:"product"`
	Discount    []Discount           `json:"discount" bson:"discount"`
	Date        time.Time            `json:"date" bson:"date"`
	Delivery    Delivery             `json:"delivery" bson:"delivery"`
	Subtotal    int                  `json:"subtotal" bson:"subtotal"`
	Status      string               `json:"status" bson:"status"`
	Status_code int                  `json:"status_code" bson:"status_code"`
	/* Status
	0. NotPayed
	1. Packed
	2. Sent
	3. Done
	4. Cenceled
	*/
	Pic_Pay string `json:"pic_pay" bson:"pic_pay"`
	To      To     `json:"to" bson:"to"`
	From    From   `json:"from" bson:"from"`
}
type To struct {
	Account Account2 `json:"account" bson:"account"`
	Name    string   `json:"name" bson:"name"`
	Number  string   `json:"number" bson:"number"`
	Address string   `json:"address" bson:"address"`
}

type From struct {
	Account Account2 `json:"account" bson:"account"`
	Name    string   `json:"name" bson:"name"`
	Number  string   `json:"number" bson:"number"`
	Address string   `json:"address" bson:"address"`
}

type TransactionModel struct{}

func (T *TransactionModel) GetStatus(status_code int) (status string) {
	switch status_code {
	case NOTPAYED:
		status = "NotPayed"
	case PACKED:
		status = "Packed"
	case SENT:
		status = "Sent"
	case DONE:
		status = "Done"
	case CENCELED:
		status = "Cenceled"
	}
	return
}

func (T *TransactionModel) Create(data forms.Transaction) (ret Transaction, err error) {
	id := uuid.New()

	// Get Account From
	data_account_from, err1 := account_model.Get(data.From.Account)
	if err1 != nil {
		err = err1
		return
	}
	// Get Account To
	data_account_to, err2 := account_model.Get(data.To.Account)
	if err2 != nil {
		err = err2
		return
	}

	// Getting the Product
	var prod []ProductTransaction
	discount_i := data.Product[0].Discount
	for i, p := range data.Product {
		var data_d Discount
		if len(data.Product) == 0 {
			prod = nil
		} else {

			_, err1 := account_model.GetDiscountUsed(data.To.Account, p.Discount)
			if err1 == nil {
				err = errors.New("discount is used on previous product")
				return
			}
			data_p, err1 := product_model.Detail(p.Product, data.From.Account)
			if err1 != nil {
				err = errors.New("Product id " + p.Product + " " + err1.Error())
				return
			}
			if i > 0 {
				if p.Discount == "" {
					data_d = Discount{}
				} else if p.Discount == discount_i {
					err = errors.New("discount used on product " + data.Product[i-1].Product)
					return
				}
			}

			if p.Discount == "" {
				data_d = Discount{}
			} else {
				data_d, err1 = discount_model.Get(p.Discount)
				if err1 != nil {
					err = errors.New("Discount Product id " + p.Discount + " " + err1.Error())
					return
				}
				err2 := account_model.AddDiscounUsed(data_account_to.Id, p.Discount)
				if err2 != nil {
					fmt.Println("log on line 96")
					err = err2
					return
				}
			}

			prod = append(prod, ProductTransaction{
				Id:      data_p.Id,
				Name:    data_p.Name,
				Qty:     p.Qty,
				Pricing: data_p.Pricing.Price * p.Qty,
				Image:   data_p.Image,
				Discount: Discount{
					Id:           data_d.Id,
					Discount:     data_d.Discount,
					DiscountCode: data_d.DiscountCode,
					Expired:      data_d.Expired,
					Image:        data_d.Image,
					Name:         data_d.Name,
				},
			})

		}
	}

	// Getting the Discount
	var dis []Discount
	for _, d := range data.Discount {
		_, err1 := account_model.GetDiscountUsed(data.To.Account, d)
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
			err2 := account_model.AddDiscounUsed(data.To.Account, d)
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
		Product:  prod,
	}

	// From
	ret.From = From{
		Account: Account2{
			Id:          data_account_from.Id,
			Name:        data_account_from.Name,
			Email:       data_account_from.Email,
			Image:       data_account_from.Image,
			Membership:  data_account_from.Membership,
			PhoneNumber: data_account_from.PhoneNumber,
			Status:      data_account_from.Status,
		},
		Address: data.From.Address,
		Name:    data.From.Name,
		Number:  data.From.Number,
	}

	// To
	ret.To = To{
		Account: Account2{
			Id:          data_account_to.Id,
			Name:        data_account_to.Name,
			Email:       data_account_to.Email,
			Image:       data_account_to.Image,
			Membership:  data_account_to.Membership,
			PhoneNumber: data_account_to.PhoneNumber,
			Status:      data_account_to.Status,
		},
		Address: data.To.Address,
		Name:    data.To.Name,
		Number:  data.To.Number,
	}

	// Proses Subtotal
	subtotal := 0
	total_discount := 0
	total := 0
	for _, dis := range prod {
		total_discount += dis.Discount.Discount
		subtotal += dis.Pricing
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

	price, _ := strconv.Atoi(data.Delivery.Price)
	ret.Delivery = Delivery{
		Courier: data.Delivery.Courier,
		Price:   price,
		Resi:    "",
		Service: data.Delivery.Service,
	}
	ret.Status = T.GetStatus(0)
	ret.Status_code = 0
	// Insert into mongo
	err = db.Collection["transaction"].Insert(ret)
	// err = db.Collection["transaction"].Insert(bson.M{
	// 	"_id":      id,
	// 	"date":     time.Now(),
	// 	"account":  data_account,
	// 	"product":  prod,
	// 	"paket":    paket,
	// 	"discount": dis,
	// 	"from": bson.M{
	// 		"account": bson.M{
	// 			"_id":         data_account_from.Id,
	// 			"name":        data_account_from.Name,
	// 			"email":       data_account_from.Email,
	// 			"phonenumber": data_account_from.PhoneNumber,
	// 			"membership":  data_account_from.Membership,
	// 			"image":       data_account_from.Image,
	// 			"status":      data_account_from.Status,
	// 		},
	// 		"name":    data.From.Name,
	// 		"number":  data.From.Number,
	// 		"address": data.From.Address,
	// 	},
	// 	"to": data.To,
	// 	"delivery": bson.M{
	// 		"courier": data.Delivery.Courier,
	// 		"service": data.Delivery.Service,
	// 		"resi":    "",
	// 		"price":   price,
	// 		"code":    data.Delivery.Code,
	// 	},
	// 	"status":      T.GetStatus(0),
	// 	"status_code": 0,
	// })
	if err != nil {
		return
	}

	return
}

func (T *TransactionModel) AddPicturePay(id_trans string, picture string) (err error) {
	path, err1 := addon.Upload("transaction", id_trans, picture)
	if err1 != nil {
		return err1
	}
	err = db.Collection["transaction"].Update(bson.M{
		"_id": id_trans,
	}, bson.M{
		"$set": bson.M{
			"pic_pay": path,
		},
	})
	return
}

func (T *TransactionModel) Get(id string) (data Transaction, err error) {
	err = db.Collection["transaction"].Find(bson.M{
		"_id": id,
	}).One(&data)
	return
}

func (T *TransactionModel) UpdateStatus(id string, status_code int) (err error) {
	if status_code == DONE {
		transaction_data, err1 := T.Get(id)
		for _, t := range transaction_data.Product {
			produck_data, _ := product_model.Get(t.Id)
			account_model.UpdatePoint(transaction_data.To.Account.Id, produck_data.Point)
		}

		if err1 != nil {
			err = errors.New("error whee getting transaction")
			return
		}
	}
	if status_code == SENT {
		transaction_data, err1 := T.Get(id)
		for _, t := range transaction_data.Product {
			produck_data, _ := product_model.Get(t.Id)

			account_model.UpdateStockProduct(transaction_data.From.Account.Id, produck_data.Id, t.Qty)
		}

		if err1 != nil {
			err = errors.New("error whee getting transaction")
			return
		}
	}
	err = db.Collection["transaction"].Update(bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			"status":      T.GetStatus(status_code),
			"status_code": status_code,
		},
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

func (T *TransactionModel) HistoyTransaction(id_account, filter, sort string, pageNo, perPage, status int) (data []Transaction, count int, err error) {
	_, err = account_model.Get(id_account)
	if err != nil {
		err = errors.New("account not found")
		return
	}
	sorting := sort
	if strings.Contains(sort, "asc") {
		sorting = strings.Replace(sort, "|asc", "", -1)
	} else if strings.Contains(sort, "desc") {
		sorting = strings.Replace(sort, "|desc", "", -1)
		sorting = "-" + sorting
	}
	regex := bson.M{"$regex": bson.RegEx{Pattern: filter, Options: "i"}}
	// pn, _ := strconv.Atoi(pageNo)
	// pp, _ := strconv.Atoi(perPage)
	err = db.Collection["transaction"].Find(bson.M{
		"status_code":      status,
		"from.account._id": id_account,
		"$or": []interface{}{
			bson.M{"product.name": regex},
		},
	}).Sort(sorting).Skip((pageNo - 1) * perPage).Limit(perPage).All(&data)
	if err != nil {
		return
	}
	count, err = db.Collection["transaction"].Find(bson.M{
		"status_code":      status,
		"from.account._id": id_account,
		"$or": []interface{}{
			bson.M{"product.name": regex},
		}}).Count()
	if err != nil {
		return
	}
	return
}
