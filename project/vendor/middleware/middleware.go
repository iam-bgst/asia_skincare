package middleware

import (
	"addon"
	"controllers"
	"expvar"
	"log"
	"time"

	exp_gin "github.com/gin-contrib/expvar"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

var (
	version = "0.0.1"
	port    = ":1998"

	productcontroll     = new(controllers.ProductControll)
	accountcontroll     = new(controllers.AccountControll)
	membershipcontroll  = new(controllers.MembershipControll)
	paketcontroll       = new(controllers.PaketControll)
	discountcontroll    = new(controllers.DiscountControll)
	transactioncontroll = new(controllers.TransactionControll)
	deliverycontroll    = new(controllers.DeliveryControll)

	// ExpVar
	counter = expvar.NewMap("counter").Init()
	last    = expvar.NewString("las_update")
)

func Middleware() {
	log.Println("Api Asia SkinCare Ready on port", port)

	var last_update struct {
		T time.Time
	}
	last_update.T = time.Now()
	last.Set(last_update.T.String())

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
	router.GET("/", HandleCounter, func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Server Asia SkinCare",
			"version": version,
		})
	})

	// Account
	account := router.Group("/account")
	{
		account.POST("/register", HandleCounter, accountcontroll.Register)
		account.GET("/checkaccount", HandleCounter, accountcontroll.CheckAccount)
		account.PUT("/update/:id", HandleCounter, accountcontroll.Update)
		account.GET("/get/:id", HandleCounter, accountcontroll.Get)
		account.PUT("/nonactive/:id", HandleCounter, accountcontroll.NonActiveAccount)
		account.PUT("/active/:id", HandleCounter, accountcontroll.ActiveAccount)
		account.GET("/list", HandleCounter, accountcontroll.ListAccount)
		account.POST("/adaddress/:id", HandleCounter, accountcontroll.AddAddress)
	}

	// Product
	product := router.Group("/product")
	{
		product.POST("/add", HandleCounter, productcontroll.Create)
		product.GET("/list", HandleCounter, productcontroll.ListByMembership)
		product.PUT("/update/:id", HandleCounter, productcontroll.Update)
		product.GET("/get/:id", HandleCounter, productcontroll.Get)
		product.PUT("/update_price/:product/:membership", HandleCounter, productcontroll.UpdatePrice)
		product.DELETE("/delete/:product", HandleCounter, productcontroll.Delete)
		product.GET("/listonagent", HandleCounter, productcontroll.ListProductOnAgent)
	}

	// Paket
	paket := router.Group("/paket")
	{
		paket.POST("/add", HandleCounter, paketcontroll.Create)
		paket.GET("/list", paketcontroll.ListByMembership)
		paket.PUT("/update/:id", HandleCounter, paketcontroll.Update)
		paket.GET("/get/:id/:idm", HandleCounter, paketcontroll.Get)
		paket.PUT("/update_product/:id", HandleCounter, paketcontroll.Updateproduct)
		paket.DELETE("/delete/:id", HandleCounter, paketcontroll.Delete)
	}

	// Discount
	discount := router.Group("/discount")
	{
		discount.POST("/add", HandleCounter, discountcontroll.Create)
		discount.GET("/list", HandleCounter, discountcontroll.List)
		discount.PUT("/update/:id", HandleCounter, discountcontroll.Update)
		discount.GET("/get/:id", HandleCounter, discountcontroll.Get)
		discount.DELETE("/delete/:id", HandleCounter, discountcontroll.Delete)
	}

	// Transaction
	transaction := router.Group("/transaction")
	{
		transaction.POST("/add", HandleCounter, transactioncontroll.Add)
		transaction.GET("/history/:account", HandleCounter, transactioncontroll.ListHistory)
		transaction.PUT("/update_status/:id", HandleCounter, transactioncontroll.UpdateStatus)
		transaction.PUT("/add_resi/:id", HandleCounter, transactioncontroll.AddResiToTransaction)
	}

	// Delivery
	delivery := router.Group("/delivery")
	{
		delivery.GET("/listcity", HandleCounter, deliverycontroll.ListCity)
		delivery.GET("/listprovince", HandleCounter, deliverycontroll.ListProvince)
		delivery.GET("/listcity_prov/:id", HandleCounter, deliverycontroll.ListCityByProvince)
		delivery.GET("/checkongkir", HandleCounter, deliverycontroll.CheckOngkir)
	}

	// Membership
	membership := router.Group("/membership")
	{
		membership.POST("/add", HandleCounter, membershipcontroll.Create)
		membership.GET("/listall", HandleCounter, membershipcontroll.ListAll)
	}

	// ExpVar
	router.GET("/debug/vars", exp_gin.Handler())
	router.Run(port)
}

func HandleCounter(c *gin.Context) {
	if c.Request.URL.Path[1:] == "" {
		counter.Add("root", 1)
	} else {
		counter.Add(c.Request.URL.Path[1:], 1)
	}

}
