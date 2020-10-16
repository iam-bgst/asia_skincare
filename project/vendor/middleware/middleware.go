package middleware

import (
	"addon"
	"controllers"
	"expvar"
	"fmt"
	"log"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
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
	metodecontroll      = new(controllers.MetodeControll)
	paymentcontroll     = new(controllers.PaymentControll)

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

	// ExpVar
	router.GET("/debug/vars", exp_gin.Handler())

	// Account
	account := router.Group("/account")
	{
		account.POST("/register", HandleCounter, accountcontroll.Register)
		account.GET("/checkaccount", HandleCounter, accountcontroll.CheckAccount)
		account.POST("/auth", HandleCounter, accountcontroll.Auth)
	}

	// Validator Jwt
	// router.Use(HandleAuth())

	{
		account.POST("/addpayment/:id", HandleCounter, accountcontroll.AddPayment)
		account.PUT("/update/:id", HandleCounter, accountcontroll.Update)
		account.GET("/get/:id", HandleCounter, accountcontroll.Get)
		account.PUT("/nonactive/:id", HandleCounter, accountcontroll.NonActiveAccount)
		account.PUT("/active/:id", HandleCounter, accountcontroll.ActiveAccount)
		account.GET("/list", HandleCounter, accountcontroll.ListAccount)
		account.POST("/adaddress/:id", HandleCounter, accountcontroll.AddAddress)
		account.POST("/qris/add/:id", HandleCounter, accountcontroll.AddQris)
	}

	// Product
	product := router.Group("/product")
	{
		product.POST("/add", HandleCounter, productcontroll.Create)
		product.GET("/list", HandleCounter, productcontroll.ListByMembership)
		product.GET("/list/recomd", HandleCounter, productcontroll.ListRecomend)
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

	// Metode
	payment := router.Group("/payment")
	{
		payment.POST("/add", HandleCounter, paymentcontroll.Create)
		payment.GET("/list", HandleCounter, paymentcontroll.List)
		payment.PUT("/update/:id", HandleCounter, paymentcontroll.Update)
		payment.GET("/get/:id", HandleCounter, paymentcontroll.Get)
		payment.DELETE("/delete/:id", HandleCounter, paymentcontroll.Delete)
		payment.GET("/listonaccount/:id", HandleCounter, paymentcontroll.ListOnAccount)
	}

	// Transaction
	transaction := router.Group("/transaction")
	{
		transaction.POST("/add", HandleCounter, transactioncontroll.Add)
		transaction.GET("/history/:account", HandleCounter, transactioncontroll.ListHistory)
		transaction.GET("/order/:account", HandleCounter, transactioncontroll.ListTransactionOnagent)
		transaction.PUT("/update_status/:id", HandleCounter, transactioncontroll.UpdateStatus)
		transaction.PUT("/add_resi/:id", HandleCounter, transactioncontroll.AddResiToTransaction)
		transaction.PUT("/add_picture/:id", HandleCounter, transactioncontroll.AddPicturePay)
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

	router.Run(port)
}

func HandleAuth() gin.HandlerFunc {
	valid := func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		Bearer := strings.Split(auth, " ")
		if Bearer[0] != "Bearer" {
			c.JSON(401, gin.H{
				"error": "Authorization do not use Bearer",
			})
			c.Abort()
		}
		token, err := jwt.Parse(Bearer[1], func(token *jwt.Token) (interface{}, error) {
			if jwt.GetSigningMethod("HS256") != token.Method {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("secret"), nil
		})

		if token != nil && err == nil {
			claims := jwt.MapClaims{}
			_, _ = jwt.ParseWithClaims(Bearer[1], claims, func(token *jwt.Token) (interface{}, error) {
				return []byte("secret"), nil
			})

		} else {
			c.JSON(401, gin.H{
				"error":   err.Error(),
				"message": "Authorization is empty",
			})
			c.Abort()
		}
	}
	return valid
}

func HandleCounter(c *gin.Context) {
	if c.Request.URL.Path[1:] == "" {
		counter.Add("root", 1)
	} else {
		counter.Add(c.Request.URL.Path[1:], 1)
	}

}
