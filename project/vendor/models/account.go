package models

import (
	"addon"
	"db"
	"errors"
	"fmt"
	"forms"
	"strconv"
	"strings"

	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2/bson"
)

type Account struct {
	Id            string     `json:"_id" bson:"_id,omitempty"`
	Name          string     `json:"name" bson:"name"`
	Email         string     `json:"email" bson:"email"`
	PhoneNumber   int        `json:"phonenumber" bson:"phonenumber"`
	Point         int        `json:"point" bson:"point"`
	Address       []Address  `json:"address" bson:"address"`
	Membership    Membership `json:"membership" bson:"membership"`
	Image         string     `json:"image" bson:"image"`
	Status        string     `json:"status" bson:"status"`
	Qris          string     `json:"qris" bson:"qris"`
	Discount_used []Discount `json:"discount_used" bson:"discount_used"`
	Product       []struct {
		Id    string `json:"_id"`
		Stock int    `json:"stock"`
	} `json:"product"`
}
type Account2 struct {
	Id          string     `json:"_id" bson:"_id,omitempty"`
	Name        string     `json:"name" bson:"name"`
	Email       string     `json:"email" bson:"email"`
	PhoneNumber int        `json:"phonenumber" bson:"phonenumber"`
	Membership  Membership `json:"membership" bson:"membership"`
	Image       string     `json:"image" bson:"image"`
	Status      string     `json:"status" bson:"status"`
}
type AccountList struct {
	Id            string `json:"_id" bson:"_id,omitempty"`
	Name          string `json:"name" bson:"name"`
	Province      string `json:"province" bson:"province"`
	Province_code int    `json:"province_id" bson:"province_code"`
	City          string `json:"city" bson:"city"`
	City_code     int    `json:"city_code" bson:"city_code"`
}

type Address struct {
	Province Province `json:"province" bson:"province"`
	City     City     `json:"city" bson:"city"`
	Detail   string   `json:"detail" bson:"detail"`
	Default  bool     `json:"default" bson:"default"`
}
type AccountTransaction struct {
	Id          string     `json:"_id" bson:"_id,omitempty"`
	Name        string     `json:"name" bson:"name"`
	Email       string     `json:"email" bson:"email"`
	PhoneNumber int        `json:"phonenumber" bson:"phonenumber"`
	Address     string     `json:"address" bson:"address"`
	Image       string     `json:"image" bson:"image"`
	Membership  Membership `json:"membership" bson:"membership"`
	Status      string     `json:"statut" bson:"status"`
}

type AccountModel struct{}

func (A *AccountModel) CheckAdmin() (found bool) {
	err := db.Collection["account"].Find(bson.M{
		"membership.code": 0,
	}).One(&bson.M{})
	if err != nil {
		return false
	} else {
		return true
	}
}

func (A *AccountModel) Create(data forms.Account) (err error) {

	id := uuid.New()
	data_membership, _ := membership_model.GetOneMembership(data.Membership)
	if data_membership.Code == STAFF && A.CheckAdmin() == false {
		return errors.New("Could not found Account Admin, you cannot create account staff while admin is nothing")
	}
	phone, _ := strconv.Atoi(data.PhoneNumber)

	path, err := addon.Upload("account", id, data.Image)
	if err != nil {
		return
	}

	prov, err1 := delivery_model.GetProvince(data.Address.Province)
	if err1 != nil {
		err = err1
		return
	}
	city, err2 := delivery_model.GetCityByProv(data.Address.Province, data.Address.City)
	if err2 != nil {
		err = err2
		return
	}

	err = db.Collection["account"].Insert(bson.M{
		"_id":         id,
		"name":        data.Name,
		"email":       data.Email,
		"phonenumber": phone,
		"membership":  data_membership,
		"point":       0,
		"image":       path,
		"status":      "active",
	})
	err = db.Collection["account"].Update(bson.M{"_id": id}, bson.M{
		"$addToSet": bson.M{
			"address": bson.M{
				"province": prov,
				"city":     city,
				"detail":   data.Address.Detail,
				"default":  true,
			},
		},
	})
	product := product_model.All()
	for _, p := range product {
		db.Collection["account"].Update(bson.M{"_id": id}, bson.M{
			"$addToSet": bson.M{
				"product": bson.M{
					"_id":   p.Id,
					"stock": 0,
				},
			},
		})
	}
	return
}

func (A *AccountModel) AddProduct(id_product string) (err error) {
	_, err = db.Collection["account"].UpdateAll(bson.M{}, bson.M{
		"$addToSet": bson.M{
			"product": bson.M{
				"_id":   id_product,
				"stock": 0,
			},
		},
	})
	return
}

func (A *AccountModel) UpdateStockOnAccount(account, product string, stock int) (err error) {
	err = db.Collection["account"].Update(bson.M{
		"_id":         account,
		"product._id": product,
	}, bson.M{
		"$set": bson.M{
			"product.$.stock": stock,
		},
	})
	return
}

func (A *AccountModel) AddQris(id_account, qris string) (err error) {
	path, err1 := addon.Upload("account/qris", id_account, qris)
	if err1 != nil {
		return err1
	}
	err = db.Collection["account"].Update(bson.M{
		"_id": account,
	}, bson.M{
		"$set": bson.M{
			"qris": path,
		},
	})
	return
}

func (A *AccountModel) AddAddress(id string, data forms.Address) (err error) {
	prov, err1 := delivery_model.GetProvince(data.Province)
	if err1 != nil {
		err = err1
		return
	}
	city, err2 := delivery_model.GetCityByProv(data.Province, data.City)
	if err2 != nil {
		err = err2
		return
	}

	err = db.Collection["account"].Update(bson.M{
		"_id": id,
	}, bson.M{
		"$addToSet": bson.M{
			"address": bson.M{
				"province": prov,
				"city":     city,
				"detail":   data.Detail,
				"default":  false,
			},
		},
	})
	return
}
func (A *AccountModel) CheckAccount(phonenumber int) (data Account, err error) {
	err = db.Collection["account"].Find(bson.M{
		"phonenumber": phonenumber,
	}).One(&data)
	return
}
func (A *AccountModel) UpdatePoint(id string, point int) (err error) {
	err = db.Collection["account"].Update(bson.M{
		"_id": id,
	}, bson.M{
		"$inc": bson.M{
			"point": point,
		},
	})
	return
}

func (A *AccountModel) UpdateStockProduct(id_account, id_product string, stock int) (err error) {
	err = db.Collection["account"].Update(bson.M{
		"_id":         id_account,
		"product._id": id_product,
	}, bson.M{
		"$inc": bson.M{
			"product.$.stock": stock - (stock * 2),
		},
	})
	return
}
func (A *AccountModel) Get(id string) (data AccountTransaction, err error) {
	err = db.Collection["account"].Find(bson.M{
		"_id": id,
	}).One(&data)
	return
}
func (A *AccountModel) GetId(id string) (data Account, err error) {
	err = db.Collection["account"].Find(bson.M{
		"_id": id,
	}).One(&data)
	return
}
func (A *AccountModel) GetByMembership(membership string, prov, city int) (available bool) {
	err := db.Collection["account"].Find(bson.M{
		"membership._id": membership,
		"$or": []interface{}{
			bson.M{"province": prov},
			bson.M{"city": city},
		},
	})
	if err == nil {
		available = true
	} else {
		available = false
	}
	return
}
func (A *AccountModel) Update(id string, data forms.Account) (err error) {
	path, err := addon.Upload("account", id, data.Image)
	if err != nil {
		return
	}
	phone, _ := strconv.Atoi(data.PhoneNumber)
	err = db.Collection["account"].Update(bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			"name":        data.Name,
			"email":       data.Email,
			"phonenumber": phone,
			"image":       path,
		},
	})
	return
}

func (A *AccountModel) NonActiveAccount(id string) (err error) {
	err = db.Collection["account"].Update(bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			"status": "nonactive",
		},
	})
	return
}

func (A *AccountModel) ActiveAccount(id string) (err error) {
	err = db.Collection["account"].Update(bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			"status": "active",
		},
	})
	return
}

func (A *AccountModel) Delete(id string) (err error) {
	err = db.Collection["account"].Remove(bson.M{
		"_id": id,
	})
	return
}

func (A *AccountModel) GetDiscountUsed(id, idd string) (data Discount, err error) {
	pipeline := []bson.M{
		{"$unwind": "$discount_used"},
		{"$match": bson.M{
			"_id":               id,
			"discount_used._id": idd,
		}},
		{"$project": bson.M{
			"_id":          "$discount_used._id",
			"name":         "$discount_used.name",
			"discount":     "$discount_used.discount",
			"discountcode": "$discount_used.discountcode",
			"image":        "$discount_used.image",
		}},
	}
	err = db.Collection["account"].Pipe(pipeline).One(&data)
	return
}

func (A *AccountModel) All() (data []Account) {
	db.Collection["account"].Find(bson.M{}).All(&data)
	return
}

func (A *AccountModel) AddDiscounUsed(id, idd string) (err error) {
	data_discount, err1 := discount_model.Get(idd)
	if err1 != nil {
		fmt.Println("log on account model line 140")
		err = err1
		return
	}
	err = db.Collection["account"].Update(bson.M{
		"_id": id,
	}, bson.M{
		"$addToSet": bson.M{
			"discount_used": bson.M{
				"_id":          data_discount.Id,
				"name":         data_discount.Name,
				"discount":     data_discount.Discount,
				"discountcode": data_discount.DiscountCode,
				"image":        data_discount.Image,
				"expired":      data_discount.Expired,
			},
		},
	})
	return
}

func (A *AccountModel) ListAccount(filter, sort string, pageNo, perPage int) (data []AccountList, count int, err error) {
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
	// pn, _ := strconv.Atoi(pageNo)
	// pp, _ := strconv.Atoi(perPage)
	regex := bson.M{"$regex": bson.RegEx{Pattern: "agen", Options: "i"}}
	pipeline := []bson.M{
		{"$match": bson.M{"membership.name": regex}},
		{"$unwind": "$address"},
		{"$match": bson.M{"address.default": true}},
		{"$project": bson.M{
			"name":          "$name",
			"province":      "$address.province.province",
			"city":          bson.M{"$concat": []string{"$address.city.city_name", " - ", "$address.city.type"}},
			"province_code": "$address.province.province_id",
			"city_code":     "$address.city.city_id",
		}},
		{"$sort": bson.M{sorting: order}},
		{"$skip": (pageNo - 1) * perPage},
		{"$limit": perPage},
	}
	err = db.Collection["account"].Pipe(pipeline).All(&data)
	count, _ = db.Collection["account"].Find(bson.M{}).Count()
	return
}
