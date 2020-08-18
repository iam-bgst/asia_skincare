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
	db.SetCollection("paket")
	db.SetCollection("discount")
}

func main() {
	// Initial Membership
	var membership = new(models.MembershipModel)
	membership.InitMembership()

	// Open Port Gin
	middleware.Middleware()

}
