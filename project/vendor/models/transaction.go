package models

import (
	"addon"
	"db"
	"errors"
	"fmt"
	"forms"
	"log"
	"strconv"
	"strings"
	"sync"
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
	Id               string               `json:"_id" bson:"_id,omitempty"`
	Transaction_code string               `json:"code" bson:"code"`
	Product          []ProductTransaction `json:"product" bson:"product"`
	Discount         []Discount           `json:"discount" bson:"discount"`
	Date             time.Time            `json:"date" bson:"date"`
	Delivery         Delivery             `json:"delivery" bson:"delivery"`
	Subtotal         int                  `json:"subtotal" bson:"subtotal"`
	Status           string               `json:"status" bson:"status"`
	Status_code      int                  `json:"status_code" bson:"status_code"`
	/* Status
	0. NotPayed
	1. Packed
	2. Sent
	3. Done
	4. Cencel
	*/
	Payment  PaymentAccount2 `json:"payment" bson:"payment"`
	Evidence Evidence        `json:"evidence" bson:"evidence"`
	To       To              `json:"to" bson:"to"`
	From     From            `json:"from" bson:"from"`
}

type Metode struct {
	Id   string `json:"_id" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	Desc string `json:"desc" bson:"desc"`
}

type Evidence struct {
	Total   string    `json:"total" bson:"total"`
	Name    string    `json:"name" bson:"name"`
	Send_by string    `json:"send_by" bson:"send_by"`
	Time    time.Time `json:"time" bson:"time"`
	Image   string    `json:"image" bson:"image"`
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

func (T *TransactionModel) Create(data forms.Transaction, wg *sync.WaitGroup) (ret Transaction, err error) {
	defer wg.Done()
	id := uuid.New()

	// Get Payment from account
	ch_payment := make(chan PaymentAccount2)
	ch_payment_err := make(chan error)
	go account_model.GetPayment(data.From.Account, data.Payment, ch_payment, ch_payment_err)

	if <-ch_payment_err != nil {
		log.Println("line 109")
		err = <-ch_payment_err
		return
	}

	// Get Account From
	data_account_from, err1 := account_model.Get(data.From.Account)
	if err1 != nil {
		log.Println("line 109")
		err = err1
		return
	}
	// Get Account To
	data_account_to, err2 := account_model.Get(data.To.Account)
	if err2 != nil {
		log.Println("line 116")
		err = err2
		return
	}

	address_to, err3 := account_model.GetAddress(data.To.Account, data.To.Address)
	if err3 != nil {
		log.Println("line 123")
		err = err3
		return
	}

	address_from, err4 := account_model.GetAddressDefault(data.From.Account)
	if err4 != nil {
		log.Println("line 130")
		err = err4
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

	ret.Payment = <-ch_payment

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
			Payment:     data_account_from.Payment,
		},
		Address: address_from.Detail,
		Name:    address_from.Name,
		Number:  address_from.Number,
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
		Address: address_to.Detail,
		Name:    address_to.Name,
		Number:  address_to.Number,
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
	ret.Transaction_code = addon.RandomCode(10, false)
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
	if err != nil {
		return
	}

	return
}

func (T *TransactionModel) TransactionOnAgent(id_account, filter, sort string, pageNo, perPage, status int) (data []Transaction, count int, err error) {
	data_account, err1 := account_model.Get(id_account)
	if err1 != nil {
		err = errors.New("account not found")
		return
	}
	var account Account
	if data_account.Membership.Code == 1 {
		account, _ = account_model.GetByCode(0)
	} else {
		account, _ = account_model.GetId(data_account.Id)
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
	fmt.Println(account.Id)
	err = db.Collection["transaction"].Find(bson.M{
		"status_code":      status,
		"from.account._id": account.Id,
		"$or": []interface{}{
			bson.M{"product.name": regex},
		},
	}).Sort(sorting).Skip((pageNo - 1) * perPage).Limit(perPage).All(&data)
	if err != nil {
		return
	}
	count, err = db.Collection["transaction"].Find(bson.M{
		"status_code":      status,
		"from.account._id": account.Id,
		"$or": []interface{}{
			bson.M{"product.name": regex},
		}}).Count()
	if err != nil {
		return
	}
	return
}

func (T *TransactionModel) AddPicturePay(id_trans string, data forms.Evidence) (err error) {
	path, err1 := addon.Upload("transaction", id_trans, data.Image)
	if err1 != nil {
		return err1
	}
	err = db.Collection["transaction"].Update(bson.M{
		"_id": id_trans,
	}, bson.M{
		"$set": bson.M{
			"evidence": bson.M{
				"image":   path,
				"send_by": data.Send_by,
				"name":    data.Name,
				"time":    data.Time,
				"total":   data.Total,
			},
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
	if status_code == SENT {
		transaction_data, err1 := T.Get(id)
		for _, t := range transaction_data.Product {

			// Getting Product
			produck_data, _ := product_model.Get(t.Id)

			// add point to account
			account_model.UpdatePoint(transaction_data.To.Account.Id, produck_data.Point)

			// add solded
			product_model.UpdateSolded(t.Id, t.Qty)

			// update stock
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
		"status_code":    status,
		"to.account._id": id_account,
		"$or": []interface{}{
			bson.M{"product.name": regex},
		},
	}).Sort(sorting).Skip((pageNo - 1) * perPage).Limit(perPage).All(&data)
	if err != nil {
		return
	}
	count, err = db.Collection["transaction"].Find(bson.M{
		"status_code":    status,
		"to.account._id": id_account,
		"$or": []interface{}{
			bson.M{"product.name": regex},
		}}).Count()
	if err != nil {
		return
	}
	return
}

func (T *TransactionModel) All() (data []Transaction) {
	db.Collection["transaction"].Find(bson.M{}).All(&data)
	return
}
