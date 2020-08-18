package controllers

import (
	"forms"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DiscountControll struct{}

func (D *DiscountControll) Create(c *gin.Context) {
	var data forms.Discount
	if c.BindJSON(&data) != nil {
		c.JSON(405, gin.H{
			"error": "error binding json",
		})
	} else {
		err := discountmodels.Create(data)
		if err != nil {
			c.JSON(405, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "Created",
				"status":  "ok",
			})
		}
	}
}

func (D *DiscountControll) Get(c *gin.Context) {
	id := c.Param("id")
	data, err := discountmodels.Get(id)
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

func (D *DiscountControll) Update(c *gin.Context) {
	id := c.Param("id")
	var data forms.Discount
	if c.BindJSON(&data) != nil {
		c.JSON(406, gin.H{
			"error": "error binding json",
		})
	} else {
		err := discountmodels.Update(id, data)
		if err != nil {
			c.JSON(405, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "updated",
				"status":  "ok",
			})
		}
	}
}

func (D *DiscountControll) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := discountmodels.Delete(id); err != nil {
		c.JSON(405, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"message": "Deleted",
			"status":  "ok",
		})
	}
}

func (D *DiscountControll) List(c *gin.Context) {
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
	data, err := discountmodels.List(sort, pageNo, perPage)
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
