package forms

type Paket struct {
	Name    string    `json:"name" bson:"name"`
	Product []string  `json:"product" bson:"product"`
	Pricing []Pricing `json:"pricing" bson:"pricing"`
	Stock   int       `json:"stock"`
	Point   int       `json:"point"`
	Image   string    `json:"image" bson:"image"`
}
