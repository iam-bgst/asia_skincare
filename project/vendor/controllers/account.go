package controllers

import (
	"forms"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AccountControll struct{}

func (A *AccountControll) Register(c *gin.Context) {
	var forms forms.Account
	if c.BindJSON(&forms) != nil {
		c.JSON(406, gin.H{"error": "Error Binding"})
	} else {
		err := accountmodels.Create(forms)
		if err != nil {
			c.JSON(406, gin.H{"error": err.Error})
		} else {
			c.JSON(200, gin.H{"message": "Registed", "status": "ok"})
		}
	}
}

func (A *AccountControll) CheckAccount(c *gin.Context) {
	phone, _ := strconv.Atoi(c.Query("phone"))
	data, err := accountmodels.CheckAccount(phone)
	if err != nil {
		c.JSON(406, gin.H{"error": err.Error})
	} else {
		c.JSON(200, gin.H{"data": data, "status": "ok"})
	}

}
