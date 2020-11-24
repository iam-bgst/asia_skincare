package swagger

import (
	"encoding/json"

	"github.com/swaggo/swag"
)

type Contact struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Email string `json:"email"`
}

type License struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Info struct {
	Description    string  `json:"description"`
	Title          string  `json:"title"`
	TermsOfService string  `json:"termsOfService"`
	Contact        Contact `json:"contact"`
	License        License `json:"license"`
	Version        string  `json:"version"`
}

type SwaggerConfig struct {
	Schemes     []string               `json:"schemes"`
	Swagger     string                 `json:"swagger"`
	Info        Info                   `json:"info"`
	Host        string                 `json:"host"`
	BasePath    string                 `json:"basePath"`
	Paths       map[string]interface{} `json:"paths"`
	Definitions map[string]interface{} `json:"definitions"`
}

var Datadoc = SwaggerConfig{
	Schemes: []string{},
	Swagger: "2.0",
	Info: Info{
		Description: "This is documentation api of Asia Skincare.",
		Title:       "Asia Skincare Api Swagger",
		Contact: Contact{
			Email: "bimagusta61@gmail.com",
			Name:  "iam.bgst",
		},
	},
	BasePath: "",
	Paths:    make(map[string]interface{}),
}

type s struct{}

func (s *s) ReadDoc() string {
	b, err := json.Marshal(Datadoc)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func init() {
	swag.Register(swag.Name, &s{})
}
