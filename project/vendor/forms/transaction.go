package forms

type Transaction struct {
	Product    []ProductTransaction `json:"product" bson:"product"`
	Paket      []string             `json:"paket" bson:"paket"`
	Discount   []string             `json:"discount" bson:"discount"`
	Membership string               `json:"membership" bson:"membership"`
	Delivery   Delivery             `json:"delivery" bson:"delivery"`
	Subtotal   int                  `json:"subtotal" bson:"subtotal"`
	To         To                   `json:"to" bson:"to"`
	From       From                 `json:"from" bson:"from"`
}

type ProductTransaction struct {
	Product  string `json:"product"`
	Discount string `json:"discount"`
}

type To struct {
	Name    string `json:"name" bson:"name"`
	Number  string `json:"number" bson:"number"`
	Address string `json:"address" bson:"address"`
}

type From struct {
	Name    string `json:"name" bson:"name"`
	Number  string `json:"number" bson:"number"`
	Address string `json:"address" bson:"address"`
}

type Delivery struct {
	Courier string `json:"courier" bson:"courier"`
	Resi    string `json:"resi" bson:"resi"`
	Price   string `json:"price" bson:"price"`
}
