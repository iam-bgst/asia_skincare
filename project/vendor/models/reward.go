package models

import (
	"addon"
	"db"
	"forms"
	"strings"
	"time"

	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2/bson"
)

type Rewards struct {
	Id         string    `json:"_id" bson:"_id,omitempty"`
	Name       string    `json:"name" bson:"name"`
	PricePoint int       `json:"pricepoint" bson:"pricepoint"`
	Desc       string    `json:"desc" bson:"desc"`
	Image      string    `json:"image" bson:"image"`
	Start      time.Time `json:"start" bson:"start"`
	End        time.Time `json:"end" bson:"end"`
}

type RewardsModels struct{}

func (R *RewardsModels) Create(data forms.Rewards) (err error) {
	id := uuid.New()
	path, err := addon.Upload("rewards", id, data.Image)
	if err != nil {
		return
	}
	err = db.Collection["rewards"].Insert(bson.M{
		"_id":        id,
		"name":       data.Name,
		"pricePoint": data.PricePoint,
		"desc":       data.Desc,
		"image":      path,
		"start":      data.Start,
		"end":        data.End,
	})
	return
}

func (R *RewardsModels) Get(id string) (data Rewards, err error) {
	err = db.Collection["rewards"].Find(bson.M{
		"_id": id,
	}).One(&data)
	return
}

func (R *RewardsModels) Update(id string, data forms.Rewards) (err error) {
	err = db.Collection["rewards"].Update(bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			"name":       data.Name,
			"pricePoint": data.PricePoint,
			"desc":       data.Desc,
			"start":      data.Start,
			"end":        data.End,
		},
	})
	return
}

func (R *RewardsModels) Delete(id string) (err error) {
	err = db.Collection["rewards"].Remove(bson.M{
		"_id": id,
	})
	return
}

func (R *RewardsModels) List(filter, sort string, pageNo, perPage int) (data []Rewards, count int, err error) {
	sorting := sort
	if strings.Contains(sort, "asc") {
		sorting = strings.Replace(sort, "|asc", "", -1)

	} else if strings.Contains(sort, "desc") {
		sorting = strings.Replace(sort, "|desc", "", -1)
		sorting = "-" + sorting

	} else {
		sorting = "date"
	}
	regex := bson.M{"$regex": bson.RegEx{Pattern: "agen", Options: "i"}}
	err = db.Collection["reward"].Find(bson.M{
		"$or": []interface{}{
			bson.M{"name": regex},
		},
	}).Sort(sorting).Skip((pageNo - 1) * perPage).Limit(perPage).All(&data)
	count, _ = db.Collection["reward"].Find(bson.M{"point.value": bson.M{"$gt": 0}}).Count()
	return
}
