package swagger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"

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

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{Schemes: []string{}}
var datadoc SwaggerConfig

type s struct{}

func (s *s) ReadDoc() string {

	b, err := json.Marshal(datadoc)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(s)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(string(b))
	if err != nil {
		return string(b)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, SwaggerInfo); err != nil {
		return string(b)
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
