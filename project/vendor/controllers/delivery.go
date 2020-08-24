package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeliveryControll struct{}

func (D *DeliveryControll) List(c *gin.Context) {
	sort, _ := strconv.Atoi(c.Query("sort"))
	pageNo, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	if pageNo == 0 {
		pageNo = 1
	}
	if perPage == 0 {
		perPage = 5
	}
	data := deliverymodels.List(sort, pageNo, perPage)
	c.JSON(200, gin.H{
		"data": data,
	})
}

func (D *DeliveryControll) CheckOngkir(c *gin.Context) {
	origin := c.Query("origin")
	destination := c.Query("destination")
	weight := c.Query("weight")
	data := deliverymodels.CheckOngkir(origin, destination, weight)
	c.JSON(200, gin.H{
		"data":   data,
		"status": "ok",
	})
}
