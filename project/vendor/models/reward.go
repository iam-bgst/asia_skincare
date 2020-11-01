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
	PricePoint int       `json:"pricepoint" bson:"pricePoint"`
	Desc       string    `json:"desc" bson:"desc"`
	Image      string    `json:"image" bson:"image"`
	Start      time.Time `json:"start" bson:"start"`
	End        time.Time `json:"end" bson:"end"`
	Archive    bool      `json:"archive" bson:"archive"`
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

func (R *RewardsModels) List(filter, sort string, pageNo, perPage int, active, archive bool) (data []Rewards, count int, err error) {
	sorting := sort
	if strings.Contains(sort, "asc") {
		sorting = strings.Replace(sort, "|asc", "", -1)

	} else if strings.Contains(sort, "desc") {
		sorting = strings.Replace(sort, "|desc", "", -1)
		sorting = "-" + sorting

	} else {
		sorting = "start"
	}
	regex := bson.M{"$regex": bson.RegEx{Pattern: filter, Options: "i"}}
	var query bson.M
	if active {
		query = bson.M{
			"start": bson.M{"$lt": time.Now()},
			"end":   bson.M{"$gt": time.Now()},
			"$or": []interface{}{
				bson.M{"name": regex},
			},
			"archive": false,
		}
	} else {
		query = bson.M{
			"$or": []interface{}{
				bson.M{"name": regex},
			},
		}
	}

	err = db.Collection["rewards"].Find(query).Sort(sorting).Skip((pageNo - 1) * perPage).Limit(perPage).All(&data)
	count, _ = db.Collection["rewards"].Find(query).Count()
	return
}

func (R *RewardsModels) Archive(id string) (err error) {
	err = db.Collection["rewards"].Update(bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			"archive": true,
		},
	})
	return
}

func (R *RewardsModels) UnArchive(id string) (err error) {
	err = db.Collection["rewards"].Update(bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			"archive": false,
		},
	})
	return
}
