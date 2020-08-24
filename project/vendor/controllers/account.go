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
		err := accountmodels.Create(data)
		if err != nil {
			c.JSON(406, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"message": "Registed", "status": "ok"})
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
