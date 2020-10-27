package models

import (
	"addon"
	"db"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pborman/uuid"
	"github.com/tidwall/gjson"
	"gopkg.in/mgo.v2/bson"
)

var city []City

type Delivery struct {
	Courier string `json:"courier" bson:"courier"`
	Service string `json:"service" bson:"service"`
	Resi    string `json:"resi" bson:"resi"`
	Price   int    `json:"price" bson:"price"`
}
type Result struct {
	Id   string  `json:"_id"`
	Name string  `json:"name"`
	Code string  `json:"code"`
	Cost []Costs `json:"cost"`
}
type Costs struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Cost     int    `json:"cost"`
	Estimate string `json:"estimate"`
}

type Province struct {
	Id          string `json:"_id" bson:"_id,omitempty"`
	Province_id int    `json:"province_id"`
	Province    string `json:"province" bson:"province"`
	City        []City `json:"city" bson:"city"`
}
type City struct {
	Id          string `json:"_id" bson:"_id,omitempty"`
	City_id     int    `json:"city_id"`
	Type        string `json:"type"`
	City_name   string `json:"city_name"`
	Postal_code int    `json:"postal_code"`
}

type DeliveryModels struct{}

var courier = []string{"jne", "pos", "tiki"}

type CheckCost struct {
	Origin      string
	Destination string
	Wieght      string
	Courier     string
}

func (D *DeliveryModels) CheckOngkirCourir(origin, destination, weight, account string) (data_result []Result) {
	var data_cost []Costs
	apiUrl := "https://api.rajaongkir.com/starter/cost?key=8c24e11e7261144361a9a5a86d30314f"

	data := url.Values{}
	data.Set("origin", origin)
	data.Set("destination", destination)
	data.Set("weight", weight)

	client := &http.Client{}
	kurir, _ := account_model.GetCourierMany(account)
	for _, s := range kurir {
		data.Set("courier", s.Code)
		r, _ := http.NewRequest("POST", apiUrl, strings.NewReader(data.Encode()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Add("key", "8c24e11e7261144361a9a5a86d30314f")

		resp, _ := client.Do(r)
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		name := gjson.Get(string(bodyBytes), "rajaongkir.results.0.name").String()

		value := gjson.Get(string(bodyBytes), "rajaongkir.results.0.costs.#.service").Array()
		for i, s := range value {
			ii := strconv.Itoa(i)
			cost := gjson.Get(string(bodyBytes), "rajaongkir.results.0.costs."+ii+".cost.0.value").Int()
			etd := gjson.Get(string(bodyBytes), "rajaongkir.results.0.costs."+ii+".cost.0.etd").String()

			data_cost = append(data_cost, Costs{
				Id:       uuid.New(),
				Name:     s.String(),
				Cost:     int(cost),
				Estimate: etd,
			})
		}
		data_result = append(data_result, Result{
			Id:   s.Id,
			Code: s.Code,
			Name: name,
			Cost: data_cost,
		})
	}
	return
}

func (D *DeliveryModels) CheckOngkir(origin, destination, weight string) (data_result []Result) {
	// var data_result []Delivery
	var data_cost []Costs
	apiUrl := "https://api.rajaongkir.com/starter/cost?key=8c24e11e7261144361a9a5a86d30314f"

	data := url.Values{}
	data.Set("origin", origin)
	data.Set("destination", destination)
	data.Set("weight", weight)

	client := &http.Client{}
	for _, s := range courier {
		data.Set("courier", s)
		r, _ := http.NewRequest("POST", apiUrl, strings.NewReader(data.Encode()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Add("key", "8c24e11e7261144361a9a5a86d30314f")

		resp, _ := client.Do(r)
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		name := gjson.Get(string(bodyBytes), "rajaongkir.results.0.name").String()

		value := gjson.Get(string(bodyBytes), "rajaongkir.results.0.costs.#.service").Array()
		for i, s := range value {
			ii := strconv.Itoa(i)
			cost := gjson.Get(string(bodyBytes), "rajaongkir.results.0.costs."+ii+".cost.0.value").Int()
			etd := gjson.Get(string(bodyBytes), "rajaongkir.results.0.costs."+ii+".cost.0.etd").String()

			data_cost = append(data_cost, Costs{
				Id:       uuid.New(),
				Name:     s.String(),
				Cost:     int(cost),
				Estimate: etd,
			})
		}
		data_result = append(data_result, Result{
			Id:   uuid.New(),
			Code: s,
			Name: name,
			Cost: data_cost,
		})
	}
	return
}

func (D *DeliveryModels) GetListCity(filter, sort string, pageNo, perPage int) (data []City, count int, err error) {
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
		sorting = "city_id"
		order = 1
	}
	regex := bson.M{"$regex": bson.RegEx{Pattern: filter, Options: "i"}}
	err = db.Collection["delivery"].Pipe([]bson.M{
		{"$project": bson.M{
			"city": "$city",
		}},
		{"$unwind": "$city"},
		{"$project": bson.M{
			"_id":         "$city._id",
			"city_id":     "$city.city_id",
			"city_name":   "$city.city_name",
			"type":        "$city.type",
			"postal_code": "$city.postal_code",
		}},
		{"$match": bson.M{
			"$or": []interface{}{
				bson.M{"city_name": regex},
			},
		}},
	}).All(&data)
	count = len(data)

	pipeline := []bson.M{
		{"$project": bson.M{
			"city": "$city",
		}},
		{"$unwind": "$city"},
		{"$project": bson.M{
			"_id":         "$city._id",
			"city_id":     "$city.city_id",
			"city_name":   "$city.city_name",
			"type":        "$city.type",
			"postal_code": "$city.postal_code",
		}},
		{"$match": bson.M{
			"$or": []interface{}{
				bson.M{"city_name": regex},
			},
		}},
		{"$sort": bson.M{sorting: order}},
		{"$skip": (pageNo - 1) * perPage},
		{"$limit": perPage},
	}

	err = db.Collection["delivery"].Pipe(pipeline).All(&data)

	return
}

func (D *DeliveryModels) GetProvince(id_province int) (data Province, err error) {
	err = db.Collection["delivery"].Pipe([]bson.M{
		{"$project": bson.M{
			"_id":         "$_id",
			"province_id": "$province_id",
			"province":    "$province",
		}},
		{"$match": bson.M{
			"province_id": id_province,
		}},
	}).One(&data)
	return
}

func (D *DeliveryModels) GetCityByProv(idprov, id_city int) (data City, err error) {
	err = db.Collection["delivery"].Pipe([]bson.M{
		{"$match": bson.M{
			"province_id": idprov,
		}},
		{"$project": bson.M{
			"city": "$city",
		}},
		{"$unwind": "$city"},
		{"$project": bson.M{
			"_id":         "$city._id",
			"city_id":     "$city.city_id",
			"city_name":   "$city.city_name",
			"type":        "$city.type",
			"postal_code": "$city.postal_code",
		}},
		{"$match": bson.M{
			"city_id": id_city,
		}},
	}).One(&data)
	return
}
func (D *DeliveryModels) GetCity(id_city int) (data City, err error) {
	err = db.Collection["delivery"].Pipe([]bson.M{
		{"$project": bson.M{
			"city": "$city",
		}},
		{"$unwind": "$city"},
		{"$project": bson.M{
			"_id":         "$city._id",
			"city_id":     "$city.city_id",
			"city_name":   "$city.city_name",
			"type":        "$city.type",
			"postal_code": "$city.postal_code",
		}},
		{"$match": bson.M{
			"city_id": id_city,
		}},
	}).One(&data)
	return
}

func (D *DeliveryModels) GetListProvince(filter, sort string, pageNo, perPage int) (data []Province, count int, err error) {
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
		sorting = "province_id"
		order = 1
	}
	regex := bson.M{"$regex": bson.RegEx{Pattern: filter, Options: "i"}}
	err = db.Collection["delivery"].Pipe([]bson.M{
		{"$project": bson.M{
			"_id":         "$_id",
			"province_id": "$province_id",
			"province":    "$province",
		}},
		{"$match": bson.M{
			"$or": []interface{}{
				bson.M{"province": regex},
			},
		}},
	}).All(&data)
	count = len(data)

	pipeline := []bson.M{
		{"$project": bson.M{
			"_id":         "$_id",
			"province_id": "$province_id",
			"province":    "$province",
		}},
		{"$match": bson.M{
			"$or": []interface{}{
				bson.M{"province": regex},
			},
		}},
		{"$sort": bson.M{sorting: order}},
		{"$skip": (pageNo - 1) * perPage},
		{"$limit": perPage},
	}

	err = db.Collection["delivery"].Pipe(pipeline).All(&data)

	return
}

func (D *DeliveryModels) GetListCityByPorvince(id_province int, filter, sort string, pageNo, perPage int) (data []City, count int, err error) {
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
		sorting = "city_id"
		order = 1
	}
	regex := bson.M{"$regex": bson.RegEx{Pattern: filter, Options: "i"}}
	err = db.Collection["delivery"].Pipe([]bson.M{
		{"$match": bson.M{
			"province_id": id_province,
		}},
		{"$project": bson.M{
			"city": "$city",
		}},
		{"$unwind": "$city"},
		{"$project": bson.M{
			"_id":         "$city._id",
			"city_id":     "$city.city_id",
			"city_name":   "$city.city_name",
			"type":        "$city.type",
			"postal_code": "$city.postal_code",
		}},
		{"$match": bson.M{
			"$or": []interface{}{
				bson.M{"city_name": regex},
			},
		}},
	}).All(&data)
	count = len(data)
	pipeline := []bson.M{
		{"$match": bson.M{
			"province_id": id_province,
		}},
		{"$project": bson.M{
			"city": "$city",
		}},
		{"$unwind": "$city"},
		{"$project": bson.M{
			"_id":         "$city._id",
			"city_id":     "$city.city_id",
			"city_name":   "$city.city_name",
			"type":        "$city.type",
			"postal_code": "$city.postal_code",
		}},
		{"$match": bson.M{
			"$or": []interface{}{
				bson.M{"city_name": regex},
			},
		}},
		{"$sort": bson.M{sorting: order}},
		{"$skip": (pageNo - 1) * perPage},
		{"$limit": perPage},
	}
	err = db.Collection["delivery"].Pipe(pipeline).All(&data)
	return
}

func (D *DeliveryModels) InitialDelivery() {
	var data []Province
	db.Collection["delivery"].Find(bson.M{}).All(&data)
	if len(data) == 0 {
		dir := addon.GetDir()
		byt, _ := ioutil.ReadFile(dir + "/vendor/config/assets.json")
		json.Unmarshal(byt, &data)
		for _, s := range data {
			id := uuid.New()
			db.Collection["delivery"].Insert(bson.M{
				"_id":         id,
				"province_id": s.Province_id,
				"province":    s.Province,
			})
			for _, c := range s.City {
				idc := uuid.New()
				db.Collection["delivery"].Update(bson.M{
					"_id": id,
				}, bson.M{
					"$addToSet": bson.M{
						"city": bson.M{
							"_id":         idc,
							"city_id":     c.City_id,
							"city_name":   c.City_name,
							"type":        c.Type,
							"postal_code": c.Postal_code,
						},
					},
				})
			}
		}
	}
}

func (D *DeliveryModels) List(sorting, pageNo, perPage int) (data []City, count int) {
	start := 0
	end := 0
	if pageNo == 1 {
		start = 0
		end = perPage
	} else if pageNo > 1 {
		start = ((perPage * pageNo) - perPage)
		end = (perPage * pageNo)
	}
	return city[start:end], len(city)
}
