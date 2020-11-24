package models

import (
	"addon"
	"db"
	"errors"
	"fmt"
	"forms"
	"strconv"
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
	code := addon.RandomCode(10, false)
	err = db.Collection["redeem"].Insert(bson.M{
		"_id":     uuid.New(),
		"code":    code,
		"reward":  data.Reward,
		"account": data.Account,
		"date":    time.Now().UTC().Add(7 * time.Hour),
		"valid":   false,
	})

	// Point Log
	d_acc, _ := account_model.GetId(data.Account)
	d_rew, _ := reward_models.Get(data.Reward)
	pointLog_model.Create(Point_log{
		Account: Account2{
			Id: data.Account,
		},
		Desc: fmt.Sprintf("Redeem #%s ", code),
		Detail: Detail{
			Type:         REDEEM,
			Code:         code,
			Point_before: d_acc.Point.Value,
			Point_after:  d_acc.Point.Value - d_rew.PricePoint,
			Point:        d_rew.PricePoint,
			Valid:        false,
		},
	})
	acc, _ := account_model.GetByCode(0)
	addon.PushNotif(acc.TokenDevice, addon.HIGH, addon.Data{
		Type:  addon.REDEEM,
		Title: "Asia SkinCare",
		Body:  fmt.Sprintf("Ada point yang ditukar dengan reward | %s ", code),
	}, "redeem|data")
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

func (R *RedeemModel) List(filter, sort string, pageNo, perPage int, valid, account string) (data []Redeem, count int, err error) {
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
	if valid != "" {
		v_valid, _ := strconv.ParseBool(valid)
		pipeline = append(pipeline,
			bson.M{"$match": bson.M{
				"valid": v_valid,
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
	if data.Valid {
		err = errors.New("Redeem Is validated")
		return
	}

	d_acc, _ := account_model.GetId(data.Account.Id)
	if d_acc.Point.Value <= 0 {
		err = errors.New("Point 0 atau tidak enough")
		addon.PushNotif(data.Account.TokenDevice, addon.HIGH, addon.Data{
			Type:  addon.REDEEM,
			Title: "Asia SkinCare",
			Body:  fmt.Sprintf("Redeem %s | Point anda tidak cukup untuk penukaran", data.Code),
		}, "redeem|redeem")
		return
	}
	err = db.Collection["redeem"].Update(bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			"valid": true,
		},
	})
	addon.PushNotif(data.Account.TokenDevice, addon.HIGH, addon.Data{
		Type:  addon.REDEEM,
		Title: "Asia SkinCare",
		Body:  fmt.Sprintf("Reward anda #%s tervalidasi oleh admin", data.Code),
	}, "redeem|redeem")
	account_model.UpdatePoint(data.Account.Id, data.Reward.PricePoint-(data.Reward.PricePoint*2))

	// Point Log
	pointLog_model.UpdateValid(data.Code, true)

	return
}
