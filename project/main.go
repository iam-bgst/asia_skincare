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
		"header",
		"point_log",
	}

	// Mongodb
	db.NewConnection()
	db.SetCollection(collection)

	// initial assets
	models.InitialAssets()

	// Service Point
	models.ServicePoint()

	// Service Transaction
	models.ServiceTransaction()

}

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @termsOfService http://swagger.io/terms/
func main() {
	// Open Port Gin
	middleware.Middleware()
}
