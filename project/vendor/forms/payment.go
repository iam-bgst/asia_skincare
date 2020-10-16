package forms

type Payment struct {
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Type   Type   `json:"type"`
	Active bool   `json:"active"`
}

type Type struct {
	Name string `json:"name"`
	Code int    `json:"code"`
}
