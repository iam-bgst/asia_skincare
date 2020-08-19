package forms

type Product struct {
	Name    string    `json:"name" bson:"name"`
	Pricing []Pricing `json:"pricing" bson:"pricing"`
	Stoct   int       `json:"stock" bson:"stock"`
	Point   int       `json:"point" bson:"point"`
	Desc    string    `json:"desc"`
}

type Pricing struct {
	Membership string `json:"membership"`
	Price      int    `json:"price"`
}
