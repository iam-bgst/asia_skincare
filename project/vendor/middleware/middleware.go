package middleware

import (
	"addon"
	"controllers"
	"expvar"
	"fmt"
	"log"
	"strings"
	"time"

	_ "swagger"

	jwt "github.com/dgrijalva/jwt-go"
	exp_gin "github.com/gin-contrib/expvar"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	rewardcontroll      = new(controllers.RewardControll)
	couriercontroll     = new(controllers.CourierControll)
	redeemcontroll      = new(controllers.RedeemControll)
	headercontroll      = new(controllers.HeaderControll)
	pointLogcontroll    = new(controllers.PointLogControll)
	router              = gin.Default()

	// ExpVar
	counter = expvar.NewMap("counter").Init()
	last    = expvar.NewString("las_update")
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	log.Println("Api Asia SkinCare Ready on port", port)

	var last_update struct {
		T time.Time
	}
	last_update.T = time.Now()
	last.Set(last_update.T.String())

	// Realese Mode

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
	// Swagger Route
	router.GET("/doc/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Assets Route
	router.Static("/public", directory+"/vendor/assets/")

	// Counter
	router.Use(HandleCounter())

	// Index
	router.GET("/", func(c *gin.Context) {
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
		account.POST("/register", accountcontroll.Register)
		account.GET("/checkaccount", accountcontroll.CheckAccount)
		account.POST("/auth", accountcontroll.Auth)
	}
	// Delivery
	delivery := router.Group("/delivery")
	{
		delivery.GET("/listcity", deliverycontroll.ListCity)
		delivery.GET("/listprovince", deliverycontroll.ListProvince)
		delivery.GET("/listcity_prov/:id", deliverycontroll.ListCityByProvince)
		delivery.GET("/checkongkir", deliverycontroll.CheckOngkir)
		delivery.GET("/track_resi", deliverycontroll.CekResi)
	}

	// Validator Jwt
	router.Use(HandleAuth())

	{
		account.GET("/account_point", accountcontroll.ListAccountPoint)
		account.PUT("/update/:id", accountcontroll.Update)
		account.GET("/get/:id", accountcontroll.Get)
		account.PUT("/nonactive/:id", accountcontroll.NonActiveAccount)
		account.PUT("/active/:id", accountcontroll.ActiveAccount)
		account.GET("/list", accountcontroll.ListAccount)
		account.POST("/adaddress/:id", accountcontroll.AddAddress)
		account.POST("/qris/add/:id", accountcontroll.AddQris)
		account.PUT("/set_token/:id", accountcontroll.SetToken)
		product_a := account.Group("/product")
		{
			product_a.POST("/add/:account", accountcontroll.AddProdcut)
			product_a.GET("/list/:account", accountcontroll.ListProduct)
			product_a.PUT("/update/:account/:product", accountcontroll.UpdateStock)
		}
		reward_a := account.Group("/reward")
		{
			reward_a.GET("/list/:reward", accountcontroll.ListAccountClaimReward)
		}
		address := account.Group("/address")
		{
			address.PUT("/update/:account/:address", accountcontroll.UpdateAddress)
			address.DELETE("/delete/:account/:address", accountcontroll.DeleteAddress)
			address.PUT("/changetodefault/:account/:address", accountcontroll.ChangeToDefault)
		}

		payment_a := account.Group("/payment")
		{
			payment_a.GET("/list/:account", accountcontroll.ListPayment)
			payment_a.POST("/add/:id", accountcontroll.AddPayment)
			payment_a.GET("/get/:account/:payment", accountcontroll.GetPayment)
			payment_a.PUT("/update/:account/:payment", accountcontroll.UpdatePayment)
			payment_a.DELETE("/delete/:account/:payment", accountcontroll.DeletePayment)
		}
		courier_a := account.Group("/courier")
		{
			courier_a.POST("/add/:account", accountcontroll.AddCourier)
			courier_a.PUT("/update/:account/:courier", accountcontroll.UpdateCourier)
			courier_a.PUT("/change/:account/:courier", accountcontroll.ActiveCourier)
			courier_a.DELETE("/delete/:account/:courier", accountcontroll.RemoveCourier)
			courier_a.GET("/list/:account", accountcontroll.ListCourier)
		}
		referral := account.Group("/referral")
		{
			referral.GET("/get/:code", accountcontroll.GetByReferralCode)
		}
	}

	// Header
	header := router.Group("/header")
	{
		header.POST("/add", headercontroll.Create)
		header.GET("/get", headercontroll.Get)
		header.PUT("/update/:id", headercontroll.Update)
	}

	// Product
	product := router.Group("/product")
	{
		product.POST("/add", productcontroll.Create)
		product.GET("/list", productcontroll.ListByMembership)
		product.GET("/list/recomd", productcontroll.ListRecomend)
		product.PUT("/update/:id", productcontroll.Update)
		product.PUT("/updateDiscount/:id", productcontroll.UpdateDiscount)
		product.GET("/get/:id", productcontroll.Get)
		product.PUT("/update_price/:product/:membership", productcontroll.UpdatePrice)
		product.DELETE("/delete/:product", productcontroll.Delete)
		product.GET("/listonagent", productcontroll.ListProductOnAgent)
		product.PUT("/archive/:product/:account", productcontroll.Archive)
		product.PUT("/unarchive/:product/:account", productcontroll.UnArchive)
	}

	// Paket
	paket := router.Group("/paket")
	{
		paket.POST("/add", paketcontroll.Create)
		paket.GET("/list", paketcontroll.ListByMembership)
		paket.PUT("/update/:id", paketcontroll.Update)
		paket.GET("/get/:id/:idm", paketcontroll.Get)
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

	// Metode
	payment := router.Group("/payment")
	{
		payment.POST("/add", paymentcontroll.Create)
		payment.GET("/list", paymentcontroll.List)
		payment.PUT("/update/:id", paymentcontroll.Update)
		payment.GET("/get/:id", paymentcontroll.Get)
		payment.DELETE("/delete/:id", paymentcontroll.Delete)
		payment.GET("/listonaccount/:id", paymentcontroll.ListOnAccount)
	}

	// Transaction
	transaction := router.Group("/transaction")
	{
		transaction.POST("/add", transactioncontroll.Add)
		transaction.GET("/history/:account", transactioncontroll.ListHistory)
		transaction.GET("/order/:account", transactioncontroll.ListTransactionOnagent)
		transaction.PUT("/update_status/:id", transactioncontroll.UpdateStatus)
		transaction.PUT("/add_resi/:id", transactioncontroll.AddResiToTransaction)
		transaction.PUT("/add_picture/:id", transactioncontroll.AddPicturePay)
	}

	// Membership
	membership := router.Group("/membership")
	{
		membership.POST("/add", membershipcontroll.Create)
		membership.GET("/listall", membershipcontroll.ListAll)
	}

	// Reward
	reward := router.Group("/reward")
	{
		reward.POST("/add", rewardcontroll.Create)
		reward.GET("/get/:id", rewardcontroll.Get)
		reward.PUT("/update/:id", rewardcontroll.Update)
		reward.GET("/list", rewardcontroll.List)
		reward.DELETE("/delete/:id", rewardcontroll.Delete)
		reward.PUT("/claim/:account/:reward", rewardcontroll.ClaimReward)
		reward.PUT("/archive/:id", rewardcontroll.Archive)
		reward.PUT("/unarchive/:id", rewardcontroll.UnArchive)
	}

	courier := router.Group("/courier")
	{
		courier.GET("/list", couriercontroll.List)
		courier.POST("/add", couriercontroll.Create)
		courier.GET("/get/:id", couriercontroll.Get)
		courier.PUT("/update/:id", couriercontroll.Update)
		courier.DELETE("/delete/:id", couriercontroll.Delete)
	}

	redeem := router.Group("/redeem")
	{
		redeem.POST("/add", redeemcontroll.Create)
		redeem.GET("/list", redeemcontroll.List)
		redeem.GET("/get/:id", redeemcontroll.Get)
		redeem.PUT("/valid/:id", redeemcontroll.Valid)
	}

	point_log := router.Group("/pointLog")
	{
		point_log.GET("/list/:account", pointLogcontroll.List)
	}
}
func Middleware() {
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

func HandleCounter() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path[1:] == "" {
			counter.Add("root", 1)
		} else {
			counter.Add(c.Request.URL.Path[1:], 1)
		}
	}
}
