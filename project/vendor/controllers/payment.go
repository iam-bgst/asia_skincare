package controllers

import (
	"forms"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentControll struct{}

func (P *PaymentControll) Create(c *gin.Context) {
	var data forms.Payment
	if c.BindJSON(&data) != nil {
		c.JSON(405, gin.H{
			"error": "error binding json",
		})
	} else {
		err := paymentmodels.Add(data)
		if err != nil {
			c.JSON(406, gin.H{
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

func (P *PaymentControll) Get(c *gin.Context) {
	id := c.Param("id")
	data, err := paymentmodels.Get(id)
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

func (P *PaymentControll) Update(c *gin.Context) {
	id := c.Param("id")
	var data forms.Payment
	if c.BindJSON(&data) != nil {
		c.JSON(406, gin.H{
			"error": "error binding json",
		})
	} else {
		err := paymentmodels.Update(id, data)
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

func (P *PaymentControll) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := paymentmodels.Delete(id); err != nil {
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

func (P *PaymentControll) List(c *gin.Context) {
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
	data, count, err := paymentmodels.List(sort, pageNo, perPage)
	lastPage := float64(count) / float64(pp)
	if pp != 0 {
		if count%pp == 0 {
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
				"per_page":     pp,
				"current_page": pn,
				"last_page":    int(lastPage),
				"next_page":    "",
				"prev_page":    "",
				"from":         ((pn * pp) - pp) + 1,
				"to":           pn * pp,
				"data":         []interface{}{},
				"status":       "Ok",
			})
			c.Abort()
		} else {
			c.JSON(200, gin.H{
				"total":        count,
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
}

func (P *PaymentControll) ListOnAccount(c *gin.Context) {
	id := c.Param("id")
	sort := c.Query("sort")
	pageNo, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	filter := c.Query("filter")

	if sort == "" {
		sort = "name"
	}
	if pageNo == 0 {
		pageNo = 1
	}
	if perPage == 0 {
		perPage = 5
	}

	data, count, err := accountmodels.ListPayment(id, filter, sort, pageNo, perPage)

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
