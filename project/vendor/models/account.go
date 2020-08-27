package models

import (
	"addon"
	"db"
	"fmt"
	"forms"
	"strconv"

	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2/bson"
)

type Account struct {
	Id            string     `json:"_id" bson:"_id,omitempty"`
	Name          string     `json:"name" bson:"name"`
	Email         string     `json:"email" bson:"email"`
	PhoneNumber   int        `json:"phonenumber" bson:"phonenumber"`
	Point         int        `json:"point" bson:"point"`
	Address       string     `json:"address" bson:"address"`
	Membership    Membership `json:"membership" bson:"membership"`
	Image         string     `json:"image" bson:"image"`
	Status        string     `json:"status" bson:"status"`
	Discount_used []Discount `json:"discount_used" bson:"discount_used"`
}
type AccountTransaction struct {
	Id          string `json:"_id" bson:"_id,omitempty"`
	Name        string `json:"name" bson:"name"`
	Email       string `json:"email" bson:"email"`
	PhoneNumber int    `json:"phonenumber" bson:"phonenumber"`
	Address     string `json:"address" bson:"address"`
	Image       string `json:"image" bson:"image"`
}

type AccountModel struct{}

func (A *AccountModel) Create(data forms.Account) (err error) {
	id := uuid.New()
	data_membership, _ := membership_model.GetOneMembership(data.Membership)
	phone, _ := strconv.Atoi(data.PhoneNumber)

	path, err := addon.Upload("account", id, data.Image)
	if err != nil {
		return
	}
	err = db.Collection["account"].Insert(bson.M{
		"_id":         id,
		"name":        data.Name,
		"email":       data.Email,
		"phonenumber": phone,
		"membership":  data_membership,
		"point":       0,
		"address":     data.Address,
		"comfirmcode": 0,
		"image":       path,
		"status":      "active",
	})
	return
}

func (A *AccountModel) CheckAccount(phonenumber int) (data Account, err error) {
	err = db.Collection["account"].Find(bson.M{
		"phonenumber": phonenumber,
	}).One(&data)
	return
}
func (A *AccountModel) UpdatePoint(id string, point int) (err error) {
	err = db.Collection["account"].Update(bson.M{
		"_id": id,
	}, bson.M{
		"$inc": bson.M{
			"point": point,
		},
	})
	return
}
func (A *AccountModel) Get(id string) (data AccountTransaction, err error) {
	err = db.Collection["account"].Find(bson.M{
		"_id": id,
	}).One(&data)
	return
}

func (A *AccountModel) Update(id string, data forms.Account) (err error) {
	phone, _ := strconv.Atoi(data.PhoneNumber)
	err = db.Collection["account"].Update(bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			"name":        data.Name,
			"email":       data.Email,
			"phonenumber": phone,
			"address":     data.Address,
		},
	})
	return
}

func (A *AccountModel) NonActiveAccount(id string) (err error) {
	err = db.Collection["account"].Update(bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			"status": "nonactive",
		},
	})
	return
}

func (A *AccountModel) ActiveAccount(id string) (err error) {
	err = db.Collection["account"].Update(bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			"status": "active",
		},
	})
	return
}

func (A *AccountModel) Delete(id string) (err error) {
	err = db.Collection["account"].Remove(bson.M{
		"_id": id,
	})
	return
}

func (A *AccountModel) GetDiscountUsed(id, idd string) (data Discount, err error) {
	pipeline := []bson.M{
		{"$unwind": "$discount_used"},
		{"$match": bson.M{
			"_id":               id,
			"discount_used._id": idd,
		}},
		{"$project": bson.M{
			"_id":          "$discount_used._id",
			"name":         "$discount_used.name",
			"discount":     "$discount_used.discount",
			"discountcode": "$discount_used.discountcode",
			"image":        "$discount_used.image",
		}},
	}
	err = db.Collection["account"].Pipe(pipeline).One(&data)
	return
}
func (A *AccountModel) AddDiscounUsed(id, idd string) (err error) {
	data_discount, err1 := discount_model.Get(idd)
	if err1 != nil {
		fmt.Println("log on account model line 140")
		err = err1
		return
	}
	err = db.Collection["account"].Update(bson.M{
		"_id": id,
	}, bson.M{
		"$addToSet": bson.M{
			"discount_used": bson.M{
				"_id":          data_discount.Id,
				"name":         data_discount.Name,
				"discount":     data_discount.Discount,
				"discountcode": data_discount.DiscountCode,
				"image":        data_discount.Image,
				"expired":      data_discount.Expired,
			},
		},
	})
	return
}
