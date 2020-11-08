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
func (P *ProductControll) ListByMembershipFix(c *gin.Context) {}

func (P *ProductControll) ListByMembership(c *gin.Context) {
	sort := c.Query("sort")
	pageNo, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	filter := c.Query("filter")
	agent := c.Query("agent")
	archive, _ := strconv.ParseBool(c.Query("archive"))
	if sort == "" {
		sort = "name"
	}
	if pageNo == 0 {
		pageNo = 1
	}
	if perPage == 0 {
		perPage = 5
	}
	var data []models.ListProducFix
	var count int
	var err error

	data, count, err = productmodels.ListProductOnAgentFix(filter, sort, pageNo, perPage, agent, archive)

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

func (P *ProductControll) ListRecomend(c *gin.Context) {
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
	var data []models.ListProducFix
	var count int
	var err error

	data, count, err = productmodels.List(filter, sort, pageNo, perPage, "")

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
				"error": err.Error(),
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
	sort := c.Query("sort")
	pageNo, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	filter := c.Query("filter")
	agent := c.Query("agent")
	archive, _ := strconv.ParseBool(c.Query("archive"))
	if sort == "" {
		sort = "name"
	}
	if pageNo == 0 {
		pageNo = 1
	}
	if perPage == 0 {
		perPage = 5
	}

	data, count, err := productmodels.ListProductOnAgent(filter, sort, pageNo, perPage, agent, archive)
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

func (P *ProductControll) Archive(c *gin.Context) {
	id := c.Param("id")
	err := productmodels.Archive(id)
	if err != nil {
		c.JSON(405, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"message": "Archived Product",
			"status":  "ok",
		})
	}
}

func (P *ProductControll) UnArchive(c *gin.Context) {
	id := c.Param("id")
	err := productmodels.UnArchive(id)
	if err != nil {
		c.JSON(405, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"message": "Unarchived Product",
			"status":  "ok",
		})
	}
}
