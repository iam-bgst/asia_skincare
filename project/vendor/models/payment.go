package models

import (
	"addon"
	"db"
	"encoding/json"
	"forms"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2/bson"
)

type Payment struct {
	Id     string `json:"_id" bson:"_id,omitempty"`
	Name   string `json:"name" bson:"name"`
	Desc   string `json:"desc" bson:"desc"`
	Type   Type   `json:"type" bson:"type"`
	Active bool   `json:"active" bson:"active"`
}

type Type struct {
	Name string `json:"name" bson:"name"`
	Code int    `json:"code" bson:"code"`
}

type PaymentModels struct{}

func (P *PaymentModels) InitialPayment() {
	var data []Payment
	db.Collection["payment"].Find(bson.M{}).All(&data)
	if len(data) == 0 {
		dir := addon.GetDir()
		byt, _ := ioutil.ReadFile(dir + "/vendor/config/type.json")
		json.Unmarshal(byt, &data)
		for _, s := range data {
			id := uuid.New()
			db.Collection["payment"].Insert(bson.M{
				"_id":    id,
				"name":   s.Name,
				"desc":   s.Desc,
				"active": s.Active,
				"type": bson.M{
					"name": s.Type.Name,
					"code": s.Type.Code,
				},
			})
		}
	}
}

func (P *PaymentModels) Add(data forms.Payment) (err error) {
	err = db.Collection["metode"].Insert(bson.M{
		"_id":  uuid.New(),
		"name": data.Name,
		"desc": data.Desc,
	})
	return
}

func (P *PaymentModels) Get(id string) (data Payment, err error) {
	err = db.Collection["metode"].Find(bson.M{
		"_id": id,
	}).One(&data)
	return
}

func (P *PaymentModels) Update(id string, data forms.Payment) (err error) {
	err = db.Collection["metode"].Update(bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			"name": data.Name,
			"desc": data.Desc,
		},
	})
	return
}

func (P *PaymentModels) Delete(id string) (err error) {
	err = db.Collection["metode"].Remove(bson.M{
		"_id": id,
	})
	return
}

func (P *PaymentModels) List(sort, pageNo, perPage string) (data []Payment, err error) {
	sorting := sort
	if strings.Contains(sort, "asc") {
		sorting = strings.Replace(sort, "|asc", "", -1)
	} else if strings.Contains(sort, "desc") {
		sorting = strings.Replace(sort, "|desc", "", -1)
		sorting = "-" + sorting
	} else {
		sorting = "date"
	}
	pn, _ := strconv.Atoi(pageNo)
	pp, _ := strconv.Atoi(perPage)
	err = db.Collection["payment"].Find(bson.M{}).Sort(sorting).Skip((pn - 1) * pp).Limit(pp).All(&data)
	return
}
