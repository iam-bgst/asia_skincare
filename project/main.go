package main

import (
	"db"
	"middleware"
	"models"
)

func init() {
	db.NewConnection()
	db.SetCollection("reseller")
	db.SetCollection("membership")
	db.SetCollection("product")
	db.SetCollection("account")
}

func main() {
	var membership = new(models.MembershipModel)
	membership.InitMembership()
	middleware.Middleware()

}
