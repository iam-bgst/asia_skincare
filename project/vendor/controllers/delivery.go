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
	data, count := deliverymodels.List(sort, pageNo, perPage)
	lastPage := float64(len(data)) / float64(perPage)
	if perPage != 0 {
		if len(data)%perPage == 0 {
			lastPage = lastPage
		} else {
			lastPage = lastPage + 1
		}
	} else {
		lastPage = float64(len(data)) / float64(5)
	}
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
