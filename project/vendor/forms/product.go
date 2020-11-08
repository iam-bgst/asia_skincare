package forms

type Product struct {
	Name     string    `json:"name" bson:"name"`
	Pricing  []Pricing `json:"pricing" bson:"pricing"`
	Weight   int       `json:"weight"`
	Netto    string    `json:"netto"`
	Stoct    int       `json:"stock" bson:"stock"`
	Point    int       `json:"point" bson:"point"`
	Desc     string    `json:"desc"`
	Image    string    `json:"image"`
	Type     int       `json:"type"`
	Discount Discount  `json:"discount"`
}

type Pricing struct {
	Membership string `json:"membership"`
	Price      int    `json:"price"`
}
