package models

import (
	"addon"
	"db"
	"forms"
	"strconv"
	"strings"
	"time"

	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2/bson"
)

type Discount struct {
	Name         string    `json:"name" bson:"name"`
	Discount     int       `json:"discount" bson:"discount"`
	DiscountCode string    `json:"discountcode" bson:"discountcode"`
	Image        string    `json:"image" bson:"image"`
	StartAt      time.Time `json:"startAt" bson:"startAt"`
	EndAt        time.Time `json:"endAt" bson:"endAt"`
	Status       bool      `json:"status" bson:"status"`
}

type DiscountModel struct{}

func (D *DiscountModel) Create(data forms.Discount) (err error) {
	var code string
	id := uuid.New()
	if data.DiscountCode == "" {
		code = addon.RandomCode(8, true)
	} else {
		code = data.DiscountCode
	}
	err = db.Collection["discount"].Insert(bson.M{
		"_id":          id,
		"name":         data.Name,
		"discount":     data.Discount,
		"discountcode": code,
		"startAt":      data.StartAt,
		"endAt":        data.EndAt,
		"product":      []interface{}{},
	})
	if err != nil {
		return
	}
	// for _, p := range data.Product {
	// 	err = db.Collection["discount"].Update(bson.M{
	// 		"_id": id,
	// 	}, bson.M{
	// 		"$addToSet": bson.M{
	// 			"product": bson.M{
	// 				"_id": p.Id,
	// 			},
	// 		},
	// 	})
	// }

	// path, err := addon.Upload("discount", id, data.Image)
	// if err != nil {
	// 	return
	// }
	// err = db.Collection["discount"].Update(bson.M{
	// 	"_id": id,
	// }, bson.M{
	// 	"$set": bson.M{
	// 		"image": path,
	// 	},
	// })
	// if err != nil {
	// 	return
	// }
	return
}

func (D *DiscountModel) Get(id string) (data Discount, err error) {
	err = db.Collection["discount"].Find(bson.M{
		"_id": id,
	}).One(&data)
	return
}

func (D *DiscountModel) Update(id string, data forms.Discount) (err error) {
	err = db.Collection["discount"].Update(bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			"name":     data.Name,
			"discount": data.Discount,
			"startAt":  data.StartAt,
			"endAt":    data.EndAt,
		},
	})
	return
}

func (D *DiscountModel) Delete(id string) (err error) {
	err = db.Collection["discount"].Remove(bson.M{
		"_id": id,
	})
	return
}

func (D *DiscountModel) List(sort, pageNo, perPage string) (data []Discount, count int, err error) {
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
	err = db.Collection["discount"].Find(bson.M{}).Sort(sorting).Skip((pn - 1) * pp).Limit(pp).All(&data)
	count, _ = db.Collection["discount"].Find(bson.M{}).Count()
	return
}
