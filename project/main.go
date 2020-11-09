package main

import (
	"db"
	"middleware"
	"models"
)

func init() {
	// Collection
	collection := []string{
		"reseller",
		"membership",
		"product",
		"payment",
		"delivery",
		"account",
		"paket",
		"discount",
		"transaction",
		"metode",
		"rewards",
		"courier",
		"redeem",
	}

	// Mongodb
	db.NewConnection()
	db.SetCollection(collection)

	// initial assets
	models.InitialAssets()

	// Service Point
	// models.ServicePoint()

	// Service Transaction
	// models.ServiceTransaction()

}

func main() {
	// Open Port Gin
	middleware.Middleware()

}
