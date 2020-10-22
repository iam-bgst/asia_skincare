package models

import (
	"addon"
	"db"
	"errors"
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
	Solded  int     `json:"solded" bson:"solded"`
	Point   int     `json:"point" bson:"point"`
	Weight  int     `json:"weight" bson:"weight"`
	Netto   string  `json:"netto": bson:"netto"`
	Image   string  `json:"image" bson:"image"`
	Desc    string  `json:"desc" bson:"desc"`
	Type    int     `json:"type" bson:"type"`
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
	Qty      int      `json:"qty" bson:"qty"`
	Image    string   `json:"image" bson:"image"`
	Pricing  int      `json:"pricing" bson:"pricing"`
	Discount Discount `json:"discount" bson:"discount"`
}

type ListProducFix struct {
	Id         string     `json:"_id" bson:"_id,omitempty"`
	Name       string     `json:"name" bson:"name"`
	Pricing    Pricing    `json:"pricing" bson:"pricing"`
	Membership Membership `json:"membership" bson:"membership"`
	Solded     int        `json:"solded" bson:"solded"`
	Stoct      int        `json:"stoct" bson:"stock"`
	Point      int        `json:"point" bson:"point"`
	Weight     int        `json:"weight" bson:"weight"`
	Netto      string     `json:"netto" bson:"netto"`
	Image      string     `json:"image" bson:"image"`
	Desc       string     `json:"desc" bson:"desc"`
	From       Address    `json:"from" bson:"from"`
	Account    string     `json:"account" bson:"account"`
	Type       int        `json:"type" bson:"type"`
}

type Pricing struct {
	Membership Membership `json:"membership" bson:"membership"`
	Price      int        `json:"price" bson:"price"`
}

type ProductModel struct{}

func (P *ProductModel) Create(data forms.Product) (err error) {
	if account_model.CheckAdmin() == false {
		return errors.New("Could not found Account Admin, admin not created")
	}
	id := uuid.New()
	path, err := addon.Upload("product", id, data.Image)
	if err != nil {
		return
	}
	err = db.Collection["product"].Insert(bson.M{
		"_id":    id,
		"name":   data.Name,
		"stock":  data.Stoct,
		"point":  data.Point,
		"desc":   data.Desc,
		"weight": data.Weight,
		"netto":  data.Netto,
		"image":  path,
		"type":   data.Type,
	})
	for _, pricing := range data.Pricing {
		data_membership, _ := membership_model.GetOneMembership(pricing.Membership)
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
	err = account_model.AddProduct(id)
	if err != nil {
		return
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

func (P *ProductModel) UpdateSolded(id string, solded int) (err error) {
	err = db.Collection["product"].Update(bson.M{
		"_id": id,
	}, bson.M{
		"$inc": bson.M{
			"solded": solded,
		},
	})
	return
}

func (P *ProductModel) Get(id string) (data Product1, err error) {
	err = db.Collection["product"].Find(bson.M{
		"_id": id,
	}).One(&data)
	return
}

func (P *ProductModel) Update(id string, data forms.Product) (err error) {
	path, err := addon.Upload("product", id, data.Image)
	if err != nil {
		return
	}
	err = db.Collection["product"].Update(bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			"name":   data.Name,
			"stock":  data.Stoct,
			"point":  data.Point,
			"weight": data.Weight,
			"netto":  data.Netto,
			"desc":   data.Desc,
			"image":  path,
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

func (P *ProductModel) All() (data []Product) {
	db.Collection["product"].Find(bson.M{}).All(&data)
	return
}

func (P *ProductModel) List(filter, sort string, pageNo, perPage int) (data []ListProducFix, count int, err error) {
	sorting := sort
	// order := 0
	if strings.Contains(sort, "asc") {
		sorting = strings.Replace(sort, "|asc", "", -1)
		// order = 1
	} else if strings.Contains(sort, "desc") {
		sorting = strings.Replace(sort, "|desc", "", -1)
		sorting = sorting
		// order = -1
	} else {
		sorting = "date"
		// order = -1
	}
	regex := bson.M{"$regex": bson.RegEx{Pattern: filter, Options: "i"}}
	pipeline := []bson.M{
		{"$match": bson.M{
			"$or": []interface{}{
				bson.M{"name": regex},
			},
		}},
		{"$sort": bson.M{"solded": -1}},
		{"$skip": (pageNo - 1) * perPage},
		{"$limit": perPage},
	}
	err = db.Collection["product"].Pipe(pipeline).All(&data)
	count, _ = db.Collection["product"].Find(bson.M{}).Count()
	return
}

func (R *ProductModel) ListByMembership(membership, filter, sort, pageNo, perPage string) (data []Product, count int, err error) {
	_, err = membership_model.GetOneMembership(membership)
	if err != nil {
		err = errors.New("membership not found")
		return
	}
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
	count, _ = db.Collection["product"].Find(bson.M{}).Count()
	return
}

func (R *ProductModel) GetByMembershipAndProvCity(membership, filter, sort, pageNo, perPage string, prov, city int) (data []Product, count int, err error) {
	ok := account_model.GetByMembership(membership, prov, city)
	if !ok {
		err = errors.New("not available agent")
		return
	}
	_, err = membership_model.GetOneMembership(membership)
	if err != nil {
		err = errors.New("membership not found")
		return
	}
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
	count, _ = db.Collection["product"].Find(bson.M{}).Count()
	return
}

func (P *ProductModel) ListProductOnAgent(filter, sort string, pageNo, perPage int, agent string) (data []ListProducFix, count int, err error) {
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
	var nin interface{}
	if agent != "" {
		nin, _ = strconv.Atoi(agent)
	} else {
		nin = bson.M{
			"$nin": []int{1, 3},
		}
	}
	pipeline := []bson.M{
		{"$match": bson.M{"membership.code": nin}},
		{"$unwind": "$product"},
		{"$unwind": "$address"},
		{"$match": bson.M{"address.default": true}},
		{"$lookup": bson.M{
			"from":         "product",
			"localField":   "product._id",
			"foreignField": "_id",
			"as":           "product_docs",
		}},
		{"$unwind": "$product_docs"},
		{"$project": bson.M{
			"_id":        "$product._id",
			"stock":      "$product.stock",
			"desc":       "$product_docs.desc",
			"from":       "$address",
			"account":    "$_id",
			"membership": "$membership",
			"name":       "$product_docs.name",
			"image":      "$product_docs.image",
			"weight":     "$product_docs.weight",
			"point":      "$product_docs.point",
			"prices":     "$product_docs.pricing",
			"netto":      "$product_docs.netto",
		}},
		{"$addFields": bson.M{
			"pricing": bson.M{
				"$arrayElemAt": []interface{}{
					bson.M{"$filter": bson.M{
						"input": "$prices",
						"as":    "pri",
						"cond": bson.M{
							"$eq": []string{"$$pri.membership._id", "$membership._id"},
						},
					},
					}, 0,
				},
			},
		}},
		{"$match": bson.M{
			"$or": []interface{}{
				bson.M{"name": regex_next},
			},
		}},
	}
	data_non_fix := []bson.M{}
	db.Collection["account"].Pipe(pipeline).All(&data_non_fix)
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
	err = db.Collection["account"].Pipe(pipeline).All(&data)
	return
}

func (P *ProductModel) Detail(id_product, id_account string) (data ListProducFix, err error) {
	pipeline := []bson.M{
		{"$match": bson.M{"_id": id_account}},
		{"$unwind": "$product"},
		{"$unwind": "$address"},
		{"$match": bson.M{"address.default": true}},
		{"$lookup": bson.M{
			"from":         "product",
			"localField":   "product._id",
			"foreignField": "_id",
			"as":           "product_docs",
		}},
		{"$unwind": "$product_docs"},
		{"$project": bson.M{
			"_id":        "$product._id",
			"stock":      "$product.stock",
			"desc":       "$product_docs.desc",
			"from":       "$address",
			"membership": "$membership",
			"name":       "$product_docs.name",
			"image":      "$product_docs.image",
			"weight":     "$product_docs.weight",
			"point":      "$product_docs.point",
			"prices":     "$product_docs.pricing",
		}},
		{"$addFields": bson.M{
			"pricing": bson.M{
				"$arrayElemAt": []interface{}{
					bson.M{"$filter": bson.M{
						"input": "$prices",
						"as":    "pri",
						"cond": bson.M{
							"$eq": []string{"$$pri.membership._id", "$membership._id"},
						},
					},
					}, 0,
				},
			},
		}},
		{"$match": bson.M{
			"_id": id_product,
		}},
	}
	err = db.Collection["account"].Pipe(pipeline).One(&data)
	return
}
