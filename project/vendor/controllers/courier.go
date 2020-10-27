package controllers

import (
	"forms"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CourierControll struct{}

func (C *CourierControll) Create(c *gin.Context) {
	var data forms.Courier
	if c.BindJSON(&data) != nil {
		c.JSON(405, gin.H{
			"error": "error binding json",
		})
	} else {
		err := couriermodels.Create(data)
		if err != nil {
			c.JSON(406, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "created courier",
				"status":  "ok",
			})
		}
	}
}

func (C *CourierControll) Get(c *gin.Context) {
	id := c.Param("id")
	data, err := couriermodels.Get(id)
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

func (C *CourierControll) Update(c *gin.Context) {
	var data forms.Courier
	if c.BindJSON(&data) != nil {
		c.JSON(405, gin.H{
			"error": "error binding json",
		})
	} else {
		id := c.Param("id")
		err := couriermodels.Update(id, data)
		if err != nil {
			c.JSON(406, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "updated courier",
				"status":  "ok",
			})
		}
	}
}

func (C *CourierControll) Delete(c *gin.Context) {
	id := c.Param("id")
	err := couriermodels.Delete(id)
	if err != nil {
		c.JSON(406, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"message": "deleted courier",
			"status":  "ok",
		})
	}
}

func (C *CourierControll) List(c *gin.Context) {
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

	data, count, err := couriermodels.List(filter, sort, pageNo, perPage)
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
