package forms

type Reseller struct {
	Id          string `json:"_id" bson:"_id,omitempty"`
	Name        string `json:"name" bson:"name"`
	Email       int    `json:"email" bson:"email"`
	PhoneNumber int    `json:"phonenumber" bson:"phonenumber"`
	Point       int    `json:"point" bson:"point"`
	Address     string `json:"address" bson:"address"`
	ConfirmCode int    `json:"confirmcode" bson:"confirmcode"`
	Membership  string `json:"membership" bson:"membership"`
	Status      string `json:"status"`
}
