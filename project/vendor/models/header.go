package models

import (
	"db"
	"forms"

	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2/bson"
)

type Header struct {
	Id    string `json:"_id" bson:"_id"`
	Title string `json:"title" bson:"title"`
}

type HeaderModel struct{}

func (H *HeaderModel) Create(data forms.Header) (err error) {
	err = db.Collection["header"].Insert(bson.M{
		"_id":   uuid.New(),
		"title": data.Title,
	})
	return
}

func (H *HeaderModel) Get() (data Header, err error) {
	err = db.Collection["header"].Find(bson.M{}).One(&data)
	return
}

func (H *HeaderModel) Update(id string, data forms.Header) (err error) {
	err = db.Collection["header"].Update(bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			"title": data.Title,
		},
	})
	return
}
