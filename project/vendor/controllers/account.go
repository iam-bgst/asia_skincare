package controllers

import (
	"forms"
	"log"
	"models"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

type AccountControll struct{}

type Payload struct {
	Id  string `json:"id"`
	Exp int    `json:"exp"`
	Jwt jwt.StandardClaims
}

func (A *AccountControll) AddPayment(c *gin.Context) {
	var data forms.AddPayment
	if c.BindJSON(&data) != nil {
		c.JSON(405, gin.H{"error": "error binding json"})
	} else {
		id := c.Param("id")
		err := accountmodels.AddPayment(id, data)
		if err != nil {
			c.JSON(406, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "added payment",
				"status":  "ok",
			})
		}
	}
}

func (A *AccountControll) SetToken(c *gin.Context) {
	var data struct {
		Token string `json:"token"`
	}
	if c.BindJSON(&data) != nil {
		c.JSON(405, gin.H{
			"error": "error binding json",
		})
	} else {
		id := c.Param("id")
		err := accountmodels.SetToken(id, data.Token)
		if err != nil {
			c.JSON(406, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "success set token",
			})
		}
	}
}

func (A *AccountControll) ListAccountPoint(c *gin.Context) {
	sort := c.Query("sort")
	pageNo, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	filter := c.Query("filter")
	if sort == "" {
		sort = "_id"
	}
	if pageNo == 0 {
		pageNo = 1
	}
	if perPage == 0 {
		perPage = 5
	}
	// pp, _ := perPage)
	// pn, _ := strconv.Atoi(pageNo)

	data, count, err := accountmodels.ListPoint(filter, sort, pageNo, perPage)
	lastPage := float64(count) / float64(perPage)
	if perPage != 0 {
		if count%perPage == 0 {
			lastPage = lastPage
		} else {
			lastPage = lastPage + 1
		}
	} else {
		lastPage = float64(count) / float64(5)
	}
	if err != nil {
		c.JSON(404, gin.H{
			"message": "terjadi kesalahan",
			"error":   err.Error(),
		})
		c.Abort()
	} else {
		if count == 0 {
			c.JSON(200, gin.H{
				"total":        count,
				"per_page":     perPage,
				"current_page": pageNo,
				"last_page":    int(lastPage),
				"next_page":    "",
				"prev_page":    "",
				"from":         ((pageNo * perPage) - perPage) + 1,
				"to":           pageNo * perPage,
				"data":         []interface{}{},
				"status":       "Ok",
			})
			c.Abort()
		} else {
			c.JSON(200, gin.H{
				"total":        count,
				"per_page":     perPage,
				"current_page": pageNo,
				"last_page":    int(lastPage),
				"next_page":    "",
				"prev_page":    "",
				"from":         ((pageNo * perPage) - perPage) + 1,
				"to":           pageNo * perPage,
				"data":         data,
				"status":       "Ok",
			})
			c.Abort()
		}
	}
}

func (A *AccountControll) GetByReferralCode(c *gin.Context) {
	code := c.Param("code")
	data, err := accountmodels.GetByReferralCode(code)
	if err != nil {
		c.JSON(405, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"data":   data,
			"status": "ok",
		})
	}
}

func (A *AccountControll) ListPayment(c *gin.Context) {
	sort := c.Query("sort")
	pageNo, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	filter := c.Query("filter")
	account := c.Param("account")
	if sort == "" {
		sort = "_id"
	}
	if pageNo == 0 {
		pageNo = 1
	}
	if perPage == 0 {
		perPage = 5
	}
	// pp, _ := perPage)
	// pn, _ := strconv.Atoi(pageNo)

	data, count, err := accountmodels.ListPayment(account, filter, sort, pageNo, perPage)
	lastPage := float64(count) / float64(perPage)
	if perPage != 0 {
		if count%perPage == 0 {
			lastPage = lastPage
		} else {
			lastPage = lastPage + 1
		}
	} else {
		lastPage = float64(count) / float64(5)
	}
	if err != nil {
		c.JSON(404, gin.H{
			"message": "terjadi kesalahan",
			"error":   err.Error(),
		})
		c.Abort()
	} else {
		if count == 0 {
			c.JSON(200, gin.H{
				"total":        count,
				"per_page":     perPage,
				"current_page": pageNo,
				"last_page":    int(lastPage),
				"next_page":    "",
				"prev_page":    "",
				"from":         ((pageNo * perPage) - perPage) + 1,
				"to":           pageNo * perPage,
				"data":         []interface{}{},
				"status":       "Ok",
			})
			c.Abort()
		} else {
			c.JSON(200, gin.H{
				"total":        count,
				"per_page":     perPage,
				"current_page": pageNo,
				"last_page":    int(lastPage),
				"next_page":    "",
				"prev_page":    "",
				"from":         ((pageNo * perPage) - perPage) + 1,
				"to":           pageNo * perPage,
				"data":         data,
				"status":       "Ok",
			})
			c.Abort()
		}
	}
}

func (A *AccountControll) UpdatePayment(c *gin.Context) {
	var data forms.AddPayment
	if c.BindJSON(&data) != nil {
		c.JSON(405, gin.H{
			"error": "Error binding json",
		})
	} else {
		id_account := c.Param("account")
		id_payment := c.Param("payment")
		err := accountmodels.UpdatePayment(id_account, id_payment, data)
		if err != nil {
			c.JSON(406, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "address updated",
				"status":  "ok",
			})
		}
	}
}

func (A *AccountControll) GetPayment(c *gin.Context) {
	account := c.Param("account")
	payment := c.Param("payment")
	data, err := accountmodels.GetPayment(account, payment)
	if err != nil {
		c.JSON(405, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"data":   data,
			"status": "ok",
		})
	}
}

func (A *AccountControll) DeletePayment(c *gin.Context) {
	id_account := c.Param("account")
	id_payment := c.Param("payment")
	err := accountmodels.DeletePayment(id_account, id_payment)
	if err != nil {
		c.JSON(406, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"message": "payment deleted",
			"status":  "ok",
		})
	}
}

func (A *AccountControll) Auth(c *gin.Context) {
	var auth struct {
		Number int `json:"number"`
	}
	if c.BindJSON(&auth) != nil {
		c.JSON(405, gin.H{
			"error": "error binding json",
		})
	} else {
		account, err := accountmodels.CheckAccount(auth.Number)
		if err != nil {
			c.JSON(406, gin.H{
				"error": err.Error(),
			})
			c.Abort()
		} else {
			payload := Payload{
				Id:  account.Id,
				Exp: int(time.Now().Add(time.Hour * 99999).Unix()),
				Jwt: jwt.StandardClaims{},
			}
			convert := jwt.NewWithClaims(jwt.SigningMethodHS256, payload.Jwt)
			token, err := convert.SignedString([]byte("secret"))
			if err != nil {
				c.JSON(406, gin.H{
					"msg":    err.Error(),
					"status": "ERROR",
				})
				c.Abort()
			} else {
				c.JSON(http.StatusOK, gin.H{
					"status":  "SUCCESS",
					"expired": time.Now().Add(time.Minute * 99999).Unix(),
					"msg":     "Sukses berhasil login",
					"token":   token,
				})
			}
		}
	}
}

func (A *AccountControll) UpdateAddress(c *gin.Context) {
	var data forms.Address
	if c.BindJSON(&data) != nil {
		c.JSON(405, gin.H{
			"error": "Error binding json",
		})
	} else {
		id_account := c.Param("account")
		id_address := c.Param("address")
		err := accountmodels.UpdateAddress(id_account, id_address, data)
		if err != nil {
			c.JSON(406, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "address updated",
				"status":  "ok",
			})
		}
	}
}

func (A *AccountControll) ChangeToDefault(c *gin.Context) {
	account := c.Param("account")
	address := c.Param("address")
	err := accountmodels.ChangeToDefault(account, address)
	if err != nil {
		c.JSON(406, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"message": "address changed to default",
			"status":  "ok",
		})
	}
}

func (A *AccountControll) DeleteAddress(c *gin.Context) {
	id_account := c.Param("account")
	id_address := c.Param("address")
	err := accountmodels.DeleteAddress(id_account, id_address)
	if err != nil {
		c.JSON(406, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"message": "address deleted",
			"status":  "ok",
		})
	}
}

func (A *AccountControll) Register(c *gin.Context) {
	log.Println("post from ip", c.ClientIP())
	var data forms.Account
	if c.BindJSON(&data) != nil {
		c.JSON(406, gin.H{
			"error": "error binding json",
		})
	} else {
		data, err := accountmodels.Create(data)
		if err != nil {
			c.JSON(406, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"message": "Registed", "data": data, "status": "ok"})
		}
	}
}

func (A *AccountControll) AddAddress(c *gin.Context) {
	id := c.Param("id")
	var data forms.Address
	if c.BindJSON(&data) != nil {
		c.JSON(400, gin.H{
			"error": "error binding json",
		})
	} else {
		err := accountmodels.AddAddress(id, data)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "added address",
				"status":  "ok",
			})
		}
	}
}

func (A *AccountControll) CheckAccount(c *gin.Context) {
	phone, _ := strconv.Atoi(c.Query("phone"))
	data, err := accountmodels.CheckAccount(phone)
	if err != nil {
		c.JSON(406, gin.H{"error": "notfound account"})
	} else {
		c.JSON(200, gin.H{"data": data, "status": "ok"})
	}
}

func (A *AccountControll) Update(c *gin.Context) {
	var forms forms.Account
	if c.BindJSON(&forms) != nil {
		c.JSON(405, gin.H{
			"error": "Error Binding json",
		})
	} else {
		id := c.Param("id")
		err := accountmodels.Update(id, forms)
		if err != nil {
			c.JSON(405, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "Updated",
				"status":  "ok",
			})
		}
	}
}

func (A *AccountControll) NonActiveAccount(c *gin.Context) {
	id := c.Param("id")
	err := accountmodels.NonActiveAccount(id)
	if err != nil {
		c.JSON(405, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"message": "NonActived",
			"status":  "ok",
		})
	}
}

func (A *AccountControll) ActiveAccount(c *gin.Context) {
	id := c.Param("id")
	err := accountmodels.ActiveAccount(id)
	if err != nil {
		c.JSON(405, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"message": "NonActived",
			"status":  "ok",
		})
	}
}

func (A *AccountControll) CheckDiscount(c *gin.Context) {
	id := c.Param("id")
	id_discount := c.Param("idd")

	data, err := accountmodels.GetDiscountUsed(id, id_discount)
	if err != nil {
		c.JSON(405, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"data":   data,
			"status": "ok",
		})
	}
}

func (A *AccountControll) ListAccount(c *gin.Context) {
	sort := c.Query("sort")
	pageNo, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	filter := c.Query("filter")
	if sort == "" {
		sort = "_id"
	}
	if pageNo == 0 {
		pageNo = 1
	}
	if perPage == 0 {
		perPage = 5
	}
	// pp, _ := perPage)
	// pn, _ := strconv.Atoi(pageNo)

	data, count, err := accountmodels.ListAccount(filter, sort, pageNo, perPage)
	lastPage := float64(count) / float64(perPage)
	if perPage != 0 {
		if count%perPage == 0 {
			lastPage = lastPage
		} else {
			lastPage = lastPage + 1
		}
	} else {
		lastPage = float64(count) / float64(5)
	}
	if err != nil {
		c.JSON(404, gin.H{
			"message": "terjadi kesalahan",
			"error":   err.Error(),
		})
		c.Abort()
	} else {
		if count == 0 {
			c.JSON(200, gin.H{
				"total":        count,
				"per_page":     perPage,
				"current_page": pageNo,
				"last_page":    int(lastPage),
				"next_page":    "",
				"prev_page":    "",
				"from":         ((pageNo * perPage) - perPage) + 1,
				"to":           pageNo * perPage,
				"data":         []interface{}{},
				"status":       "Ok",
			})
			c.Abort()
		} else {
			c.JSON(200, gin.H{
				"total":        count,
				"per_page":     perPage,
				"current_page": pageNo,
				"last_page":    int(lastPage),
				"next_page":    "",
				"prev_page":    "",
				"from":         ((pageNo * perPage) - perPage) + 1,
				"to":           pageNo * perPage,
				"data":         data,
				"status":       "Ok",
			})
			c.Abort()
		}
	}
}

func (A *AccountControll) AddQris(c *gin.Context) {
	id := c.Param("id")
	var data struct {
		Name  string `json:"name"`
		NMID  string `json:"nmid"`
		Image string `json:"image"`
	}
	if c.BindJSON(&data) != nil {
		c.JSON(405, gin.H{
			"error": "error bind json",
		})
	} else {
		err := accountmodels.AddQris(id, data.Image, data.Name, data.NMID)
		if err != nil {
			c.JSON(406, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "added qris",
				"status":  "ok",
			})
		}
	}
}

func (A *AccountControll) Get(c *gin.Context) {
	id := c.Param("id")
	data, err := accountmodels.GetId(id)
	if err != nil {
		c.JSON(405, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"data":   data,
			"status": "ok",
		})
	}
}
func (A *AccountControll) AddProdcut(c *gin.Context) {
	var data struct {
		Id    string `json:"_id"`
		Stock int    `json:"stock"`
	}
	if c.BindJSON(&data) != nil {
		c.JSON(405, gin.H{
			"error": "error binding json",
		})
	} else {
		account := c.Param("account")
		err := accountmodels.AddProduct(account, data.Id, data.Stock)
		if err != nil {
			c.JSON(406, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "Product added",
				"status":  "ok",
			})
		}
	}
}

func (A *AccountControll) ListProduct(c *gin.Context) {
	sort := c.Query("sort")
	pageNo, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	filter := c.Query("filter")
	account := c.Param("account")
	if sort == "" {
		sort = "_id"
	}
	if pageNo == 0 {
		pageNo = 1
	}
	if perPage == 0 {
		perPage = 5
	}
	// pp, _ := perPage)
	// pn, _ := strconv.Atoi(pageNo)

	data, count, err := productmodels.List(filter, sort, pageNo, perPage, account)
	lastPage := float64(count) / float64(perPage)
	if perPage != 0 {
		if count%perPage == 0 {
			lastPage = lastPage
		} else {
			lastPage = lastPage + 1
		}
	} else {
		lastPage = float64(count) / float64(5)
	}
	if err != nil {
		c.JSON(404, gin.H{
			"message": "terjadi kesalahan",
			"error":   err.Error(),
		})
		c.Abort()
	} else {
		if count == 0 {
			c.JSON(200, gin.H{
				"total":        count,
				"per_page":     perPage,
				"current_page": pageNo,
				"last_page":    int(lastPage),
				"next_page":    "",
				"prev_page":    "",
				"from":         ((pageNo * perPage) - perPage) + 1,
				"to":           pageNo * perPage,
				"data":         []interface{}{},
				"status":       "Ok",
			})
			c.Abort()
		} else {
			c.JSON(200, gin.H{
				"total":        count,
				"per_page":     perPage,
				"current_page": pageNo,
				"last_page":    int(lastPage),
				"next_page":    "",
				"prev_page":    "",
				"from":         ((pageNo * perPage) - perPage) + 1,
				"to":           pageNo * perPage,
				"data":         data,
				"status":       "Ok",
			})
			c.Abort()
		}
	}
}

func (A *AccountControll) UpdateStock(c *gin.Context) {
	account := c.Param("account")
	product := c.Param("product")
	var data struct {
		Stock int `json:"stock"`
	}
	if c.BindJSON(&data) != nil {
		c.JSON(405, gin.H{
			"error": "error binding json",
		})
	} else {
		err := accountmodels.UpdateStockOnAccount(account, product, data.Stock)
		if err != nil {
			c.JSON(406, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "stock updated",
				"status":  "ok",
			})
		}
	}

}

func (A *AccountControll) ListAccountClaimReward(c *gin.Context) {
	sort := c.Query("sort")
	pageNo, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	filter := c.Query("filter")
	reward := c.Param("reward")

	if sort == "" {
		sort = "name"
	}
	if pageNo == 0 {
		pageNo = 1
	}
	if perPage == 0 {
		perPage = 5
	}
	var data []models.Account
	var count int
	var err error

	data, count, err = accountmodels.ListAccountRewardClaim(filter, sort, reward, pageNo, perPage)

	lastPage := float64(count) / float64(perPage)
	if perPage != 0 {
		if len(data)%perPage == 0 {
			lastPage = lastPage
		} else {
			lastPage = lastPage + 1
		}
	} else {
		lastPage = float64(count) / float64(5)
	}
	if err != nil {
		c.JSON(404, gin.H{
			"message": "terjadi kesalahan",
			"error":   err.Error(),
		})
		c.Abort()
	} else {
		if count == 0 {
			c.JSON(200, gin.H{
				"total":        count,
				"per_page":     perPage,
				"current_page": pageNo,
				"last_page":    int(lastPage),
				"next_page":    "",
				"prev_page":    "",
				"from":         ((pageNo * perPage) - perPage) + 1,
				"to":           pageNo * perPage,
				"data":         []interface{}{},
				"status":       "Ok",
			})
			c.Abort()
		} else {
			c.JSON(200, gin.H{
				"total":        count,
				"per_page":     perPage,
				"current_page": pageNo,
				"last_page":    int(lastPage),
				"next_page":    "",
				"prev_page":    "",
				"from":         ((pageNo * perPage) - perPage) + 1,
				"to":           pageNo * perPage,
				"data":         data,
				"status":       "Ok",
			})
			c.Abort()
		}

	}
}

// :COURIER
func (A *AccountControll) AddCourier(c *gin.Context) {
	var data forms.AddCourier
	if c.BindJSON(&data) != nil {
		c.JSON(405, gin.H{
			"error": "error binding json",
		})
	} else {
		account := c.Param("account")
		err := accountmodels.AddCourier(account, data)
		if err != nil {
			c.JSON(406, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "courier added",
				"status":  "ok",
			})
		}
	}
}

func (A *AccountControll) UpdateCourier(c *gin.Context) {
	var data forms.AddCourier
	if c.BindJSON(&data) != nil {
		c.JSON(405, gin.H{
			"error": "error binding json",
		})
	} else {
		account := c.Param("account")
		courier := c.Param("courier")
		err := accountmodels.UpdateCourier(account, courier, data)
		if err != nil {
			c.JSON(406, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "courier updated",
				"status":  "ok",
			})
		}
	}
}

func (A *AccountControll) ActiveCourier(c *gin.Context) {
	account := c.Param("account")
	courier := c.Param("courier")
	active, _ := strconv.ParseBool(c.Query("active"))

	err := accountmodels.ActiveCourier(account, courier, active)
	if err != nil {
		c.JSON(406, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"message": "courier updated",
			"status":  "ok",
		})
	}

}

func (A *AccountControll) RemoveCourier(c *gin.Context) {
	account := c.Param("account")
	courier := c.Param("courier")

	err := accountmodels.RemoveCourier(account, courier)
	if err != nil {
		c.JSON(406, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"message": "courier deleted",
			"status":  "ok",
		})
	}
}

func (A *AccountControll) ListCourier(c *gin.Context) {
	sort := c.Query("sort")
	pageNo, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	filter := c.Query("filter")
	account := c.Param("account")
	active, _ := strconv.ParseBool(c.Query("active"))
	if sort == "" {
		sort = "_id"
	}
	if pageNo == 0 {
		pageNo = 1
	}
	if perPage == 0 {
		perPage = 5
	}
	// pp, _ := perPage)
	// pn, _ := strconv.Atoi(pageNo)

	data, count, err := accountmodels.ListCourier(account, filter, sort, pageNo, perPage, active)
	lastPage := float64(count) / float64(perPage)
	if perPage != 0 {
		if count%perPage == 0 {
			lastPage = lastPage
		} else {
			lastPage = lastPage + 1
		}
	} else {
		lastPage = float64(count) / float64(5)
	}
	if err != nil {
		c.JSON(404, gin.H{
			"message": "terjadi kesalahan",
			"error":   err.Error(),
		})
		c.Abort()
	} else {
		if count == 0 {
			c.JSON(200, gin.H{
				"total":        count,
				"per_page":     perPage,
				"current_page": pageNo,
				"last_page":    int(lastPage),
				"next_page":    "",
				"prev_page":    "",
				"from":         ((pageNo * perPage) - perPage) + 1,
				"to":           pageNo * perPage,
				"data":         []interface{}{},
				"status":       "Ok",
			})
			c.Abort()
		} else {
			c.JSON(200, gin.H{
				"total":        count,
				"per_page":     perPage,
				"current_page": pageNo,
				"last_page":    int(lastPage),
				"next_page":    "",
				"prev_page":    "",
				"from":         ((pageNo * perPage) - perPage) + 1,
				"to":           pageNo * perPage,
				"data":         data,
				"status":       "Ok",
			})
			c.Abort()
		}
	}
}
