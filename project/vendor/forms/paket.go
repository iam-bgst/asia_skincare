package forms

type Paket struct {
	Id      string    `json:"_id" bson:"_id,omitempty"`
	Name    string    `json:"name" bson:"name"`
	Product []string  `json:"product" bson:"product"`
	Pricing []Pricing `json:"pricing" bson:"pricing"`
	Stock   int       `json:"stock"`
	Image   string    `json:"image" bson:"image"`
}
