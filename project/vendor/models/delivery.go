package models

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/parnurzeal/gorequest"
	"github.com/tidwall/gjson"
)

var city []City

type Delivery struct {
	Courier string `json:"courier" bson:"courier"`
	Service string `json:"service" bson:"service"`
	Resi    string `json:"resi" bson:"resi"`
	Price   int    `json:"price" bson:"price"`
}
type Result struct {
	Name string  `json:"name"`
	Code string  `json:"code"`
	Cost []Costs `json:"cost"`
}
type Costs struct {
	Name     string `json:"name"`
	Cost     int    `json:"cost"`
	Estimate string `json:"estimate"`
}

type City struct {
	City_id     string `json:"city_id"`
	Province_id string `json:"province_id"`
	Province    string `json:"province"`
	Type        string `json:"type"`
	City_name   string `json:"city_name"`
	Posta_code  string `json:"postal_code"`
}

type DeliveryModels struct{}

var courier = []string{"jne", "pos", "tiki"}

type CheckCost struct {
	Origin      string
	Destination string
	Wieght      string
	Courier     string
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
		// fmt.Println(resp.Status)
		name := gjson.Get(string(bodyBytes), "rajaongkir.results.0.name").String()

		// data_result = append(data_result, Delivery{
		// 	Code: s,
		// 	Name: name,
		// })
		// fmt.Println(string(bodyBytes))
		value := gjson.Get(string(bodyBytes), "rajaongkir.results.0.costs.#.service").Array()
		for i, s := range value {
			ii := strconv.Itoa(i)
			cost := gjson.Get(string(bodyBytes), "rajaongkir.results.0.costs."+ii+".cost.0.value").Int()
			etd := gjson.Get(string(bodyBytes), "rajaongkir.results.0.costs."+ii+".cost.0.etd").String()

			data_cost = append(data_cost, Costs{
				Name:     s.String(),
				Cost:     int(cost),
				Estimate: etd,
			})
		}
		data_result = append(data_result, Result{
			Code: s,
			Name: name,
			Cost: data_cost,
		})
	}
	// fmt.Println(data_result)
	return
}

func (D *DeliveryModels) GetListCity() {
	req := gorequest.New()
	_, body, _ := req.Get("https://api.rajaongkir.com/starter/city?key=8c24e11e7261144361a9a5a86d30314f").End()
	value := gjson.Get(body, "rajaongkir.results").String()
	// fmt.Println(value)
	json.Unmarshal([]byte(value), &city)
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
