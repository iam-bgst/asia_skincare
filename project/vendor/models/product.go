package models

import (
	"db"
	"forms"
	"strconv"
	"strings"

	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2/bson"
)

type Product struct {
	Id      string  `json:"_id" bson:"_id,omitempty"`
	Name    string  `json:"name" bson:"name"`
	Pricing Pricing `json:"pricing" bson:"pricing"`
	Stoct   int     `json:"stoct" bson:"stock"`
	Point   int     `json:"point" bson:"point"`
	Image   string  `json:"image" bson:"image"`
}
type Product1 struct {
	Id      string    `json:"_id" bson:"_id,omitempty"`
	Name    string    `json:"name" bson:"name"`
	Pricing []Pricing `json:"pricing" bson:"pricing"`
	Stoct   int       `json:"stoct" bson:"stock"`
	Point   int       `json:"point" bson:"point"`
	Image   string    `json:"image" bson:"image"`
}
type Product2 struct {
	Id    string `json:"_id" bson:"_id,omitempty"`
	Name  string `json:"name" bson:"name"`
	Image string `json:"image" bson:"image"`
}

type ProductMembership struct {
	Id      string  `json:"_id" bson:"_id,omitempty"`
	Name    string  `json:"name" bson:"name"`
	Image   string  `json:"image" bson:"image"`
	Pricing Pricing `json:"pricing" bson:"pricing"`
}
type ProductTransaction struct {
	Id       string   `json:"_id" bson:"_id,omitempty"`
	Name     string   `json:"name" bson:"name"`
	Image    string   `json:"image" bson:"image"`
	Pricing  int      `json:"pricing" bson:"pricing"`
	Discount Discount `json:"discount" bson:"discount"`
}
type Pricing struct {
	Membership Membership `json:"membership" bson:"membership"`
	Price      int        `json:"price" bson:"price"`
}

type ProductModel struct{}

func (P *ProductModel) Create(data forms.Product) (err error) {
	id := uuid.New()
	err = db.Collection["product"].Insert(bson.M{
		"_id":   id,
		"name":  data.Name,
		"stock": data.Stoct,
		"point": data.Point,
	})
	for _, pricing := range data.Pricing {
		data_membership := membership_model.GetOneMembership(pricing.Membership)
		err = db.Collection["product"].Update(bson.M{
			"_id": id,
		}, bson.M{
			"$addToSet": bson.M{
				"pricing": bson.M{
					"membership": data_membership,
					"price":      pricing.Price,
				},
			},
		})
	}
	return
}

func (P *ProductModel) GetByMembership(id, idm string) (data ProductMembership, err error) {
	pipeline := []bson.M{
		{"$unwind": "$pricing"},
		{"$match": bson.M{
			"_id":                    id,
			"pricing.membership._id": idm,
		}},
	}
	err = db.Collection["product"].Pipe(pipeline).One(&data)
	return
}

func (P *ProductModel) Get(id string) (data Product1, err error) {
	err = db.Collection["product"].Find(bson.M{
		"_id": id,
	}).One(&data)
	return
}

func (P *ProductModel) Update(id string, data forms.Product) (err error) {
	err = db.Collection["product"].Update(bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			"name":  data.Name,
			"stock": data.Stoct,
			"point": data.Point,
		},
	})
	return
}

func (P *ProductModel) UpdatePriceByMembership(id_product, id_membership string, price int) (err error) {
	err = db.Collection["product"].Update(bson.M{
		"_id":                    id_product,
		"pricing.membership._id": id_membership,
	}, bson.M{
		"$set": bson.M{
			"pricing.$.price": price,
		},
	})
	return
}

func (R *ProductModel) Delete(id string) (err error) {
	err = db.Collection["product"].Remove(bson.M{
		"_id": id,
	})
	return
}

func (R *ProductModel) ListByMembership(membership, filter, sort, pageNo, perPage string) (data []Product, err error) {
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
	pn, _ := strconv.Atoi(pageNo)
	pp, _ := strconv.Atoi(perPage)
	regex := bson.M{"$regex": bson.RegEx{Pattern: filter, Options: "i"}}
	pipeline := []bson.M{
		{"$match": bson.M{
			"$or": []interface{}{
				bson.M{"name": regex},
			},
		}},
		{"$unwind": "$pricing"},
		{"$match": bson.M{"pricing.membership._id": membership}},
		{"$sort": bson.M{sorting: order}},
		{"$skip": (pn - 1) * pp},
		{"$limit": pp},
	}
	err = db.Collection["product"].Pipe(pipeline).All(&data)
	return
}
