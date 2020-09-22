package controllers

import (
	"forms"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AccountControll struct{}

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

func (A *AccountControll) AddQris(c *gin.Context) {
	id := c.Param("id")
	var data struct {
		Image string `json:"image"`
	}
	if c.BindJSON(&data) != nil {
		c.JSON(405, gin.H{
			"error": "error bind json",
		})
	} else {
		err := accountmodels.AddQris(id, data.Image)
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
