package models

import (
	"addon"
	"db"
	"fmt"
	"forms"
	"mime/multipart"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2/bson"
)

type Account struct {
	Id          string     `json:"_id" bson:"_id,omitempty"`
	Name        string     `json:"name" bson:"name"`
	Email       string     `json:"email" bson:"email"`
	PhoneNumber int        `json:"phonenumber" bson:"phonenumber"`
	Point       int        `json:"point" bson:"point"`
	Address     string     `json:"address" bson:"address"`
	ConfirmCode int        `json:"confirmcode" bson:"confirmcode"`
	Membership  Membership `json:"membership" bson:"membership"`
	Image       string     `json:"image" bson:"image"`
	Status      string     `json:"status" bson:"status"`
}

type AccountModel struct{}

func (A *AccountModel) Create(data forms.Account, file *multipart.FileHeader, c *gin.Context) (err error) {
	id := uuid.New()
	data_membership := membership_model.GetOneMembership(data.Membership)
	phone, _ := strconv.Atoi(data.PhoneNumber)
	fmt.Println(data)
	fmt.Println(file)

	path, _ := addon.Upload("account", id, file, c)
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

func (A *AccountModel) Get(id string) (data Account, err error) {
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
