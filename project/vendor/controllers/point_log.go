package controllers

import (
	"models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PointLogControll struct{}

func (P *PointLogControll) List(c *gin.Context) {
	sort := c.Query("sort")
	pageNo, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	filter := c.Query("filter")
	account := c.Param("account")
	if sort == "" {
		sort = "name"
	}
	if pageNo == 0 {
		pageNo = 1
	}
	if perPage == 0 {
		perPage = 5
	}
	var data []models.Point_log
	var count int
	var err error

	data, count, err = pointLogmodels.List(account, filter, sort, pageNo, perPage)

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
