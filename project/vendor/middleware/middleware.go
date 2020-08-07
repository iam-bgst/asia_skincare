package middleware

import (
	"controllers"

	"github.com/gin-gonic/gin"
)

var (
	productcontroll = new(controllers.ProductControll)
	accountcontroll = new(controllers.AccountControll)
)

func Middleware() {
	router := gin.Default()

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
