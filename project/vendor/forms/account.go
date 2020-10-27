package forms

type Account struct {
	Name        string  `json:"name" bson:"name"`
	Email       string  `json:"email" bson:"email"`
	PhoneNumber string  `json:"phonenumber" bson:"phonenumber"`
	Address     Address `json:"address" bson:"address"`
	Membership  string  `json:"membership" bson:"membership"`
	Image       string  `json:"image"`
}

type Address struct {
	Name     string `json:"name"`
	ZipCode  string `json:"zipcode"`
	Number   string `json:"number"`
	Province int    `json:"province"`
	City     int    `json:"city"`
	Detail   string `json:"detail"`
}

type AddPayment struct {
	Id     string `json:"_id"`
	An     string `json:"an"`
	Number string `json:"number"`
}

type AddCourier struct {
	Id string `json:"_id"`
}
