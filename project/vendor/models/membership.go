package models

import (
	"db"
	"strconv"

	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2/bson"
)

type Membership struct {
	Id   string `json:"_id" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	Code int    `json:"code" bson:"code"`
}

type MembershipModel struct{}

var Memship []string

const (
	CENTER = iota
	STAFF
	AGENT
	RESELLER
)

func (MS *MembershipModel) GetOneMembership(id string) (data Membership, err error) {
	err = db.Collection["membership"].Find(bson.M{
		"_id": id,
	}).One(&data)
	return
}

func (MS *MembershipModel) InitMembership() (err error) {
	Memship = []string{"Admin", "Staff", "Distributor", "Agent", "Reseller"}
	// fmt.Println(db.Collection["membership"])
	// db.SetCollection("membership")
	var data []Membership
	db.Collection["membership"].Find(bson.M{}).All(&data)
	if len(data) == 0 {
		for code, member := range Memship {
			id := uuid.New()
			db.Collection["membership"].Insert(bson.M{
				"_id":  id,
				"name": member,
				"code": code,
			})
		}
		err = nil
	}
	return
}

func (MS *MembershipModel) ListAll(ne, code string) (data []Membership, err error) {
	neq := []int{}
	nn, _ := strconv.Atoi(ne)
	if code != "" {

		if ne != "" {
			neq = []int{1, nn}
		} else {
			neq = []int{1}
		}
		err = db.Collection["membership"].Find(bson.M{
			"code": bson.M{
				"$nin": neq,
			},
		}).Sort("code").All(&data)
		return
	}
	if account_model.CheckAdmin() {
		if nn > 0 {
			err = db.Collection["membership"].Find(bson.M{"code": bson.M{
				"$nin": []int{nn, 0},
			}}).Sort("code").All(&data)
		} else {
			err = db.Collection["membership"].Find(bson.M{"code": bson.M{
				"$ne": 0,
			}}).Sort("code").All(&data)
		}
	} else {
		if nn > 0 {
			err = db.Collection["membership"].Find(bson.M{"code": bson.M{
				"$ne": ne,
			}}).Sort("code").All(&data)
		} else {
			err = db.Collection["membership"].Find(bson.M{}).Sort("code").All(&data)
		}
	}
	return
}

func (MS *MembershipModel) Add(name string) (err error) {
	id := uuid.New()
	err = db.Collection["membership"].Insert(bson.M{
		"_id":  id,
		"name": name,
	})
	return
}

func (MS *MembershipModel) GetAdminAccount(data Membership, err error) {
	err = db.Collection["membership"].Find(bson.M{
		"code": 0,
	}).One(&data)
	return
}
