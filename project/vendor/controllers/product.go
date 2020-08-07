package controllers

import (
	"forms"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductControll struct{}

func (P *ProductControll) Create(c *gin.Context) {
	var data forms.Product
	if c.BindJSON(&data) != nil {
		c.JSON(406, gin.H{"error": "binding json error"})
	} else {
		err := productmodels.Create(data)
		if err != nil {
			c.JSON(400, gin.H{"error": err})
		} else {
			c.JSON(200, gin.H{"status": "ok", "message": "Created"})
		}
	}
}

func (P *ProductControll) ListByMembership(c *gin.Context) {
	membership := c.Query("membership")
	sort := c.Query("sort")
	pageNo := c.Query("page")
	perPage := c.Query("per_page")
	if sort == "" {
		sort = "id"
	}
	if pageNo == "" {
		pageNo = "1"
	}
	if perPage == "" {
		perPage = "5"
	}
	pp, _ := strconv.Atoi(perPage)
	pn, _ := strconv.Atoi(pageNo)

	data, err := productmodels.ListByMembership(membership, sort, pageNo, perPage)
	lastPage := float64(len(data)) / float64(pp)
	if pp != 0 {
		if len(data)%pp == 0 {
			lastPage = lastPage
		} else {
			lastPage = lastPage + 1
		}
	} else {
		lastPage = float64(len(data)) / float64(5)
	}
	if err != nil {
		c.JSON(404, gin.H{
			"message": "terjadi kesalahan",
			"error":   err.Error(),
		})
		c.Abort()
	} else {
		c.JSON(200, gin.H{
			"total":        len(data),
			"per_page":     pp,
			"current_page": pn,
			"last_page":    int(lastPage),
			"next_page":    "",
			"prev_page":    "",
			"from":         ((pn * pp) - pp) + 1,
			"to":           pn * pp,
			"data":         data,
			"status":       "Ok",
		})
		c.Abort()
	}
}
