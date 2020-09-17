package models

import (
	"db"

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
	Memship = []string{"Admin", "Staff", "Agent", "Reseller"}
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

func (MS *MembershipModel) ListAll(ne int) (data []Membership, err error) {
	if account_model.CheckAdmin() {
		if ne > 0 {
			err = db.Collection["membership"].Find(bson.M{"code": bson.M{
				"$nin": []int{ne, 0},
			}}).All(&data)
		} else {
			err = db.Collection["membership"].Find(bson.M{"code": bson.M{
				"$ne": 0,
			}}).All(&data)
		}
	} else {
		if ne > 0 {
			err = db.Collection["membership"].Find(bson.M{"code": bson.M{
				"$ne": ne,
			}}).All(&data)
		} else {
			err = db.Collection["membership"].Find(bson.M{}).All(&data)
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
