package controllers

import (
	"forms"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AccountControll struct{}

func (A *AccountControll) Register(c *gin.Context) {
	var data forms.Account
	file, _, _ := c.Request.FormFile("image")
	data.Address = c.PostForm("address")
	data.Email = c.PostForm("email")
	data.Membership = c.PostForm("membership")
	data.Name = c.PostForm("name")
	data.PhoneNumber = c.PostForm("phonenumber")

	err := accountmodels.Create(data, file)
	if err != nil {
		c.JSON(406, gin.H{"error": err.Error})
	} else {
		c.JSON(200, gin.H{"message": "Registed", "status": "ok"})
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
