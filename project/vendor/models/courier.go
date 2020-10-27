package models

import (
	"db"
	"forms"
	"strings"

	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2/bson"
)

type Courier struct {
	Id   string `json:"_id" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	Code string `json:"code" bson:"code"`
	Desc string `json:"desc" bson:"desc"`
}

type CourierModel struct{}

func (C *CourierModel) InitialCourier() {
	var data []Courier
	db.Collection["courier"].Find(bson.M{}).All(&data)
	if len(data) == 0 {
		for _, c := range courier {
			C.Create(forms.Courier{
				Name: c,
				Code: strings.ToLower(c),
				Desc: "",
			})
		}
	} else {
		acc := account_model.All()
		kurir := C.All()
		go func() {
			for _, ac := range acc {
				if len(ac.Courier) == 0 {
					for _, co := range kurir {
						account_model.AddCourier(ac.Id, forms.AddCourier{
							Id: co.Id,
						})
					}
				}
			}
		}()
	}
}

func (C *CourierModel) Create(data forms.Courier) (err error) {
	id := uuid.New()
	err = db.Collection["courier"].Insert(bson.M{
		"_id":  id,
		"name": data.Name,
		"code": data.Code,
		"desc": data.Desc,
	})
	acc := account_model.All()
	go func() {
		for _, ac := range acc {
			account_model.AddCourier(ac.Id, forms.AddCourier{
				Id: id,
			})
		}
	}()

	return
}

func (C *CourierModel) Get(id string) (data Courier, err error) {
	err = db.Collection["courier"].Find(bson.M{
		"_id": id,
	}).One(&data)
	return
}

func (C *CourierModel) GetCode(code string) (data Courier, err error) {
	err = db.Collection["courier"].Find(bson.M{
		"code": code,
	}).One(&data)
	return
}
func (C *CourierModel) Update(id string, data forms.Courier) (err error) {
	err = db.Collection["courier"].Update(bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			"name": data.Name,
			"code": data.Code,
			"desc": data.Desc,
		},
	})
	return
}

func (C *CourierModel) Delete(id string) (err error) {
	err = db.Collection["courier"].Remove(bson.M{
		"_id": id,
	})
	return
}

func (C *CourierModel) List(filter, sort string, pageNo, perPage int) (data []Courier, count int, err error) {
	sorting := sort
	if strings.Contains(sort, "asc") {
		sorting = strings.Replace(sort, "|asc", "", -1)
	} else if strings.Contains(sort, "desc") {
		sorting = "-" + strings.Replace(sort, "|desc", "", -1)
	} else {
		sorting = "name"
	}
	err = db.Collection["courier"].Find(bson.M{}).Sort(sorting).Skip((pageNo - 1) * perPage).Limit(perPage).All(&data)
	count, _ = db.Collection["courier"].Find(bson.M{}).Count()
	return
}

func (C *CourierModel) All() (data []Courier) {
	db.Collection["courier"].Find(bson.M{}).All(&data)
	return
}
