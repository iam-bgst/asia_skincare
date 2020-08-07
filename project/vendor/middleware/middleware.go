package middleware

import (
	"controllers"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

var (
	productcontroll = new(controllers.ProductControll)
	accountcontroll = new(controllers.AccountControll)
)

func Middleware() {
	router := gin.Default()
	// Cors
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type, Device",
		ExposedHeaders:  "",
		Credentials:     true,
		ValidateHeaders: false,
	}))

	account := router.Group("/account")
	{
		account.POST("/register", accountcontroll.Register)
		account.GET("/checkaccount", accountcontroll.CheckAccount)
	}

	product := router.Group("/product")
	{
		product.POST("/add", productcontroll.Create)
		product.GET("/list", productcontroll.ListByMembership)
	}

	router.Run(":1998")
}
