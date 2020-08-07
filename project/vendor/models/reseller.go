package models

import (
	"db"
	"forms"

	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2/bson"
)

type Reseller struct {
	Id          string     `json:"_id" bson:"_id,omitempty"`
	Name        string     `json:"name" bson:"name"`
	Email       int        `json:"email" bson:"email"`
	PhoneNumber int        `json:"phonenumber" bson:"phonenumber"`
	Point       int        `json:"point" bson:"point"`
	Address     string     `json:"address" bson:"address"`
	ConfirmCode int        `json:"confirmcode" bson:"confirmcode"`
	Membership  Membership `json:"membership" bson:"membership"`
	Image       string     `json:"image" bson:"image"`
	Status      string     `json:"status" bson:"status"`
}

type ResellerModel struct{}

func (R *ResellerModel) Create(data forms.Reseller) (err error) {
	id := uuid.New()
	data_membership := membership_model.GetOneMembership(data.Membership)
	err = db.Collection["reseller"].Insert(bson.M{
		"_id":         id,
		"name":        data.Name,
		"email":       data.Email,
		"phonenumber": data.PhoneNumber,
		"membership":  data_membership,
		"point":       0,
		"address":     data.Address,
		"comfirmcode": 0,
		"status":      "active",
	})
	return
}

func (R *ResellerModel) Get(id string) (data Reseller, err error) {
	err = db.Collection["reseller"].Find(bson.M{
		"_id": id,
	}).One(&data)
	return
}

func (R *ResellerModel) Update(id string, data forms.Reseller) (err error) {
	err = db.Collection["reseller"].Update(bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			"name":        data.Name,
			"email":       data.Email,
			"phonenumber": data.PhoneNumber,
			"point":       data.Point,
			"address":     data.Address,
		},
	})
	return
}

func (R *ResellerModel) NonActiveAccount(id string) (err error) {
	err = db.Collection["reseller"].Update(bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			"status": "nonactive",
		},
	})
	return
}

func (R *ResellerModel) Delete(id string) (err error) {
	err = db.Collection["reseller"].Remove(bson.M{
		"_id": id,
	})
	return
}
