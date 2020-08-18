package middleware

import (
	"addon"
	"controllers"
	"log"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

var (
	version = "0.0.1"
	port    = ":1998"

	productcontroll    = new(controllers.ProductControll)
	accountcontroll    = new(controllers.AccountControll)
	membershipcontroll = new(controllers.MembershipControll)
	paketcontroll      = new(controllers.PaketControll)
	discountcontroll   = new(controllers.DiscountControll)
)

func Middleware() {
	log.Println("Api Asia SkinCare Ready on port", port)
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	directory := addon.GetDir()

	// Cors
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Access-Control-Allow-Methods,Origin, Authorization, Content-Type, X-Request-With",
		ExposedHeaders:  "",
		Credentials:     true,
		ValidateHeaders: false,
	}))
	router.Static("/public", directory+"/vendor/assets/")

	// Index
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Server Asia SkinCare",
			"version": version,
		})
	})

	// Account
	account := router.Group("/account")
	{
		account.POST("/register", accountcontroll.Register)
		account.GET("/checkaccount", accountcontroll.CheckAccount)
		account.PUT("/update/:id", accountcontroll.Update)
		account.PUT("/nonactive/:id", accountcontroll.NonActiveAccount)
		account.PUT("/active/:id", accountcontroll.ActiveAccount)
	}

	// Product
	product := router.Group("/product")
	{
		product.POST("/add", productcontroll.Create)
		product.GET("/list", productcontroll.ListByMembership)
		product.PUT("/update/:id", productcontroll.Update)
		product.GET("/get/:id", productcontroll.Get)
		product.PUT("/update_price/:product/:membership", productcontroll.UpdatePrice)
		product.DELETE("/delete/:product", productcontroll.Delete)
	}

	// Paket
	paket := router.Group("/paket")
	{
		paket.POST("/add", paketcontroll.Create)
		paket.GET("/list", paketcontroll.ListByMembership)
		paket.PUT("/update/:id", paketcontroll.Update)
		paket.GET("/get/:id", paketcontroll.Get)
		paket.PUT("/update_product/:id", paketcontroll.Updateproduct)
		paket.DELETE("/delete/:id", paketcontroll.Delete)
	}

	// Discount
	discount := router.Group("/discount")
	{
		discount.POST("/add", discountcontroll.Create)
		discount.GET("/list", discountcontroll.List)
		discount.PUT("/update/:id", discountcontroll.Update)
		discount.GET("/get/:id", discountcontroll.Get)
		discount.DELETE("/delete/:id", discountcontroll.Delete)
	}

	// Membership
	membership := router.Group("/membership")
	{
		membership.POST("/add", membershipcontroll.Create)
		membership.GET("/listall", membershipcontroll.ListAll)
	}
	router.Run(port)
}
