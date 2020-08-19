package controllers

import (
	"forms"

	"github.com/gin-gonic/gin"
)

type TransactionControll struct{}

func (T *TransactionControll) Add(c *gin.Context) {
	var data forms.Transaction
	if c.BindJSON(&data) != nil {
		c.JSON(405, gin.H{
			"error": "error binding json",
		})
	} else {
		da, err := transactionmodels.Create(data)
		if err != nil {
			c.JSON(405, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "Transaction Created",
				"data":    da,
			})
		}
	}
}
