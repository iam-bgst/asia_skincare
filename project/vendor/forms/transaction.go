package forms

type Transaction struct {
	Account    string               `json:"account"`
	Product    []ProductTransaction `json:"product" bson:"product"`
	Paket      []PaketTransaction   `json:"paket" bson:"paket"`
	Discount   []string             `json:"discount" bson:"discount"`
	Membership string               `json:"membership" bson:"membership"`
	Delivery   Delivery             `json:"delivery" bson:"delivery"`
	Subtotal   int                  `json:"subtotal" bson:"subtotal"`
	To         To                   `json:"to" bson:"to"`
	From       From                 `json:"from" bson:"from"`
}

type PaketTransaction struct {
	Paket string `json:"paket"`
	Qty   int    `json:"qty"`
}

type ProductTransaction struct {
	Product  string `json:"product"`
	Qty      int    `json:"qty"`
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
	Courier string `json:"courier"`
	Service string `json:"service"`
	Resi    string `json:"resi"`
	Price   string `json:"price"`
	Code    string `json:"code"`
}
