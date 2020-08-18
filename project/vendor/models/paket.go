package models

import (
	"addon"
	"db"
	"forms"
	"strconv"
	"strings"

	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2/bson"
)

type Paket struct {
	Id      string     `json:"_id" bson:"_id,omitempty"`
	Name    string     `json:"name" bson:"name"`
	Product []Product2 `json:"product" bson:"product"`
	Pricing Pricing    `json:"pricing" bson:"pricing"`
	Stock   int        `json:"stock" bson:"stock"`
	Point   int        `json:"point" bson:"point"`
	Image   string     `json:"image" bson:"image"`
}

type PaketModel struct{}

func (P *PaketModel) Create(data forms.Paket) (err error) {
	id := uuid.New()
	path, err := addon.Upload("paket", id, data.Image)
	if err != nil {
		return
	}
	err = db.Collection["paket"].Insert(bson.M{
		"_id":   id,
		"name":  data.Name,
		"image": path,
		"stock": data.Stock,
		"point": data.Point,
	})
	for _, product := range data.Product {
		data_product, _ := product_model.Get(product)
		err = db.Collection["paket"].Update(bson.M{
			"_id": id,
		}, bson.M{
			"$addToSet": bson.M{
				"product": bson.M{
					"_id":   data_product.Id,
					"name":  data_product.Name,
					"image": data_product.Image,
				},
			},
		})
	}
	for _, pricing := range data.Pricing {
		data_membership := membership_model.GetOneMembership(pricing.Membership)
		err = db.Collection["paket"].Update(bson.M{
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

func (P *PaketModel) Get(id string) (data Paket, err error) {
	err = db.Collection["paket"].Find(bson.M{
		"_id": id,
	}).One(&data)
	return
}

func (P *PaketModel) Update(id string, data forms.Paket) (err error) {
	err = db.Collection["paket"].Update(bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			"name":  data.Name,
			"stock": data.Stock,
			"point": data.Point,
		},
	})
	return
}

func (P *PaketModel) UpdateProduct(id string, id_product []string) (err error) {
	err = db.Collection["paket"].Update(bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			"product": []interface{}{},
		},
	})
	for _, product := range id_product {
		data_product, _ := product_model.Get(product)
		err = db.Collection["paket"].Update(bson.M{
			"_id": id,
		}, bson.M{
			"$addToSet": bson.M{
				"product": bson.M{
					"_id":  data_product.Id,
					"name": data_product.Name,
				},
			},
		})
		if err != nil {
			return
		}
	}

	return
}

func (P *PaketModel) Delete(id string) (err error) {
	err = db.Collection["paket"].Remove(bson.M{
		"_id": id,
	})
	return
}

func (P *PaketModel) ListByMembership(membership, filter, sort, pageNo, perPage string) (data []Paket, err error) {
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
		// {"$or": []interface{}{
		// 	bson.M{"name": regex},
		// }},
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
	err = db.Collection["paket"].Pipe(pipeline).All(&data)
	return
}
