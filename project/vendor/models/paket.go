package models

import (
	"db"
	"forms"

	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2/bson"
)

type Paket struct {
	Id      string    `json:"_id" bson:"_id,omitempty"`
	Name    string    `json:"name" bson:"name"`
	Product []Product `json:"product" bson:"product"`
	Pricing []Pricing `json:"pricing" bson:"pricing"`
	Stock   int       `json:"stock" bson:"stock"`
	Image   string    `json:"image" bson:"image"`
}

type PaketModel struct{}

func (P *PaketModel) Create(data forms.Paket) (err error) {
	id := uuid.New()
	err = db.Collection["paket"].Insert(bson.M{
		"_id":   id,
		"name":  data.Name,
		"image": "",
		"stock": data.Stock,
	})
	for _, product := range data.Product {
		data_product, _ := product_model.Get(product)
		err = db.Collection["paket"].Update(bson.M{
			"_id": id,
		}, bson.M{
			"$addToSet": bson.M{
				"product": bson.M{
					"_id":     data_product.Id,
					"name":    data_product.Name,
					"stock":   0,
					"point":   0,
					"pricing": []interface{}{},
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
