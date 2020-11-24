package models

import (
	"db"
	"strings"
	"time"

	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2/bson"
)

type Point_log struct {
	Id      string    `json:"_id" bson:"_id,omitempty"`
	Date    time.Time `json:"date" bson:"date"`
	Desc    string    `json:"desc" bson:"desc"`
	Account Account2  `json:"account" bson:"account"`
	Detail  Detail    `json:"detail" bson:"detail"`
}

/*
Point Log
0. Transaksi
1. Redeem
*/
const (
	TRANSACTION = iota
	REDEEM
)

type Detail struct {
	Type         int    `json:"type"`
	Code         string `json:"code"`
	Point_after  int    `json:"point_after"`
	Point_before int    `json:"point_before"`
	Point        int    `json:"point"`
	Valid        bool   `json:"valid"`
}

type Point_log_Model struct{}

func (P *Point_log_Model) Create(data Point_log) (err error) {
	if data.Id == "" {
		data.Id = uuid.New()
	}
	err = db.Collection["point_log"].Insert(bson.M{
		"_id":     data.Id,
		"date":    time.Now(),
		"desc":    data.Desc,
		"account": data.Account.Id,
		"detail":  data.Detail,
	})
	return
}

func (P *Point_log_Model) List(account, filter, sort string, pageNo, perPage int) (data []Point_log, count int, err error) {
	sorting := sort
	order := 0
	if strings.Contains(sort, "asc") {
		sorting = strings.Replace(sort, "|asc", "", -1)
		order = 1
	} else if strings.Contains(sort, "desc") {
		sorting = strings.Replace(sort, "|desc", "", -1)
		sorting = sorting
		order = -1
	} else {
		sorting = "date"
		order = -1
	}
	pipeline := []bson.M{}
	pipeline = []bson.M{
		{"$match": bson.M{
			"account": account,
		}},
		{"$lookup": bson.M{
			"from":         "account",
			"localField":   "account",
			"foreignField": "_id",
			"as":           "account",
		}},
		{"$unwind": "$account"},
	}
	data_non_fix := []bson.M{}

	db.Collection["point_log"].Pipe(pipeline).All(&data_non_fix)

	count = len(data_non_fix)
	err = db.Collection["point_log"].Pipe(pipeline).All(&data)
	pipeline = append(pipeline,
		bson.M{"$sort": bson.M{sorting: order}},
	)
	pipeline = append(pipeline,
		bson.M{"$skip": (pageNo - 1) * perPage},
	)
	pipeline = append(pipeline,
		bson.M{"$limit": perPage},
	)
	err = db.Collection["point_log"].Pipe(pipeline).All(&data)
	return
}

func (P *Point_log_Model) UpdateValid(code string, valid bool) (err error) {
	err = db.Collection["point_log"].Update(bson.M{
		"detail.code": code,
	}, bson.M{
		"$set": bson.M{
			"detail.valid": valid,
		},
	})
	return
}
