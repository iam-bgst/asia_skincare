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

type Redeem struct {
	Id      string             `json:"_id" bson:"_id"`
	Account AccountTransaction `json:"account" bson:"account"`
	Code    string             `json:"code" bson:"code"`
	Reward  Rewards            `json:"reward" bson:"reward"`
	Date    time.Time          `json:"date" bson:"date"`
	Valid   bool               `json:"valid" bson:"valid"`
}

type RedeemModel struct{}

func (R *RedeemModel) Create(data forms.Redeem) (data_return Redeem, err error) {
	_, err0 := reward_models.Get(data.Reward)
	if err0 != nil {
		return (Redeem{}), err0
	}

	_, err1 := account_model.Get(data.Account)
	if err1 != nil {
		return (Redeem{}), err1
	}

	err = db.Collection["redeem"].Insert(bson.M{
		"_id":     uuid.New(),
		"code":    addon.RandomCode(10, false),
		"reward":  data.Reward,
		"account": data.Account,
		"date":    time.Now(),
		"valid":   false,
	})
	return
}
func (R *RedeemModel) Get(id string) (data Redeem, err error) {
	pipeline := []bson.M{
		{"$match": bson.M{
			"_id": id,
		}},
		{"$lookup": bson.M{
			"from":         "rewards",
			"localField":   "reward",
			"foreignField": "_id",
			"as":           "reward",
		}},
		{"$unwind": "$reward"},
		{"$lookup": bson.M{
			"from":         "account",
			"localField":   "account",
			"foreignField": "_id",
			"as":           "account",
		}},
		{"$unwind": "$account"},
	}
	err = db.Collection["redeem"].Pipe(pipeline).One(&data)
	return
}

func (R *RedeemModel) List(filter, sort string, pageNo, perPage int, valid bool, account string) (data []Redeem, count int, err error) {
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

	regex_next := bson.M{"$regex": bson.RegEx{Pattern: filter, Options: "i"}}
	pipeline := []bson.M{
		{"$match": bson.M{
			"valid": valid,
		}},
		{"$lookup": bson.M{
			"from":         "rewards",
			"localField":   "reward",
			"foreignField": "_id",
			"as":           "reward",
		}},
		{"$unwind": "$reward"},
		{"$lookup": bson.M{
			"from":         "account",
			"localField":   "account",
			"foreignField": "_id",
			"as":           "account",
		}},
		{"$unwind": "$account"},
		{"$match": bson.M{
			"$or": []interface{}{
				bson.M{"reward.name": regex_next},
			},
		}},
	}
	if account != "" {
		pipeline = append(pipeline,
			bson.M{"$match": bson.M{
				"account._id": account,
			}},
		)
	}

	data_non_fix := []bson.M{}
	db.Collection["redeem"].Pipe(pipeline).All(&data_non_fix)
	count = len(data_non_fix)
	pipeline = append(pipeline,
		bson.M{"$sort": bson.M{sorting: order}},
	)
	pipeline = append(pipeline,
		bson.M{"$skip": (pageNo - 1) * perPage},
	)
	pipeline = append(pipeline,
		bson.M{"$limit": perPage},
	)
	err = db.Collection["redeem"].Pipe(pipeline).All(&data)
	return
}

func (R *RedeemModel) Valid(id string) (err error) {
	data, _ := R.Get(id)

	err = db.Collection["redeem"].Update(bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			"valid": true,
		},
	})
	account_model.UpdatePoint(data.Account.Id, data.Reward.PricePoint-(data.Reward.PricePoint*2))
	return
}
