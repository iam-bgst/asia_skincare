package controllers

import (
	"forms"
	"models"
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
	filter := c.Query("filter")

	prov, _ := strconv.Atoi(c.Query("prov"))
	city, _ := strconv.Atoi(c.Query("city"))

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
	var data []models.Product
	var count int
	var err error

	if prov == 0 || city == 0 {
		data, count, err = productmodels.ListByMembership(membership, filter, sort, pageNo, perPage)
	} else {
		data, count, err = productmodels.GetByMembershipAndProvCity(membership, filter, sort, pageNo, perPage, prov, city)
	}
	lastPage := float64(count) / float64(pp)
	if pp != 0 {
		if len(data)%pp == 0 {
			lastPage = lastPage
		} else {
			lastPage = lastPage + 1
		}
	} else {
		lastPage = float64(count) / float64(5)
	}
	if membership == "" {
		c.JSON(400, gin.H{
			"error": "membership not set",
		})
		c.Abort()
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

func (P *ProductControll) Get(c *gin.Context) {
	id := c.Param("id")
	data, err := productmodels.Get(id)
	if err != nil {
		c.JSON(406, gin.H{
			"error": err,
		})
	} else {
		c.JSON(200, gin.H{
			"status": "ok",
			"data":   data,
		})
	}
}

func (P *ProductControll) Update(c *gin.Context) {
	var forms forms.Product
	if c.BindJSON(&forms) != nil {
		c.JSON(405, gin.H{"error": "Error while binding json"})
	} else {
		id := c.Param("id")
		err := productmodels.Update(id, forms)
		if err != nil {
			c.JSON(406, gin.H{
				"error": "error update product",
			})
		} else {
			c.JSON(200, gin.H{
				"message": "success update product",
				"status":  "ok",
			})
		}
	}
}

func (P *ProductControll) UpdatePrice(c *gin.Context) {
	id_product := c.Param("product")
	id_membership := c.Param("membership")

	price, _ := strconv.Atoi(c.Query("price"))

	err := productmodels.UpdatePriceByMembership(id_product, id_membership, price)
	if err != nil {
		c.JSON(406, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"message": "Update Price",
			"status":  "ok",
		})
	}
}

func (P *ProductControll) Delete(c *gin.Context) {
	id_product := c.Param("product")
	err := productmodels.Delete(id_product)
	if err != nil {
		c.JSON(405, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"message": "Deleted Product",
			"status":  "ok",
		})
	}
}

func (P *ProductControll) ListProductOnAgent(c *gin.Context) {
	id_account_agent := c.Param("id_account_agent")
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

	data, count, err := productmodels.ListProductOnAgent(id_account_agent, filter, sort, pageNo, perPage)
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
