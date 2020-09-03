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
	db.SetCollection("transaction")
	db.SetCollection("delivery")

	// initial assets
	models.InitialAssets()

}

func main() {
	// Open Port Gin
	middleware.Middleware()

}
