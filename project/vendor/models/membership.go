package models

import (
	"db"

	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2/bson"
)

type Membership struct {
	Id   string `json:"_id" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
}

type MembershipModel struct{}

var Memship []string

func (MS *MembershipModel) GetOneMembership(id string) (data Membership) {
	db.Collection["membership"].Find(bson.M{
		"_id": id,
	}).One(&data)
	return
}

func (MS *MembershipModel) InitMembership() (err error) {
	Memship = []string{"Admin", "Reseller", "Reseller Agen", "Staf"}
	// fmt.Println(db.Collection["membership"])
	// db.SetCollection("membership")
	var data []Membership
	db.Collection["membership"].Find(bson.M{}).All(&data)
	if len(data) == 0 {
		for _, member := range Memship {
			id := uuid.New()
			db.Collection["membership"].Insert(bson.M{
				"_id":  id,
				"name": member,
			})
		}
		err = nil
	}
	return
}

func (MS *MembershipModel) ListAll() (data []MembershipModel, err error) {
	err = db.Collection["membership"].Find(bson.M{}).All(&data)
	return
}

func (MS *MembershipModel) Add(name string) (err error) {
	id := uuid.New()
	err = db.Collection[""].Insert(bson.M{
		"_id":  id,
		"name": name,
	})
	return
}
