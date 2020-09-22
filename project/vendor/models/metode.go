package models

import (
	"db"
	"forms"
	"strconv"
	"strings"

	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2/bson"
)

type MetodeModels struct{}

func (M *MetodeModels) Add(data forms.Metode) (err error) {
	err = db.Collection["metode"].Insert(bson.M{
		"_id":  uuid.New(),
		"name": data.Name,
		"desc": data.Desc,
	})
	return
}

func (M *MetodeModels) Get(id string) (data Metode, err error) {
	err = db.Collection["metode"].Find(bson.M{
		"_id": id,
	}).One(&data)
	return
}

func (M *MetodeModels) Update(id string, data forms.Metode) (err error) {
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

func (M *MetodeModels) Delete(id string) (err error) {
	err = db.Collection["metode"].Remove(bson.M{
		"_id": id,
	})
	return
}

func (M *MetodeModels) List(sort, pageNo, perPage string) (data []Metode, err error) {
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
	err = db.Collection["metode"].Find(bson.M{}).Sort(sorting).Skip((pn - 1) * pp).Limit(pp).All(&data)
	return
}
