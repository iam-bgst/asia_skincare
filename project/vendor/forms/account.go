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
	Province int    `json:"province"`
	City     int    `json:"city"`
	Detail   string `json:"detail"`
}
