package main

import (
	"db"
	"middleware"
	"models"
)

func init() {
	// Collection
	collection := []string{"reseller", "membership", "product", "delivery", "account", "paket", "discount", "transaction", "metode", "reward"}

	// Mongodb
	db.NewConnection()
	db.SetCollection(collection)

	// initial assets
	models.InitialAssets()

	// Service Point
	models.ServicePoint()

}

func main() {
	// Open Port Gin
	middleware.Middleware()

}
