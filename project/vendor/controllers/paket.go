package controllers

import (
	"forms"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaketControll struct{}

func (P *PaketControll) Create(c *gin.Context) {
	var data forms.Paket
	if c.BindJSON(&data) != nil {
		c.JSON(406, gin.H{"error": "error binding json"})
	} else {
		err := paketmodels.Create(data)
		if err != nil {
			c.JSON(406, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "created",
				"status":  "ok",
			})
		}
	}
	return
}

func (P *PaketControll) Get(c *gin.Context) {
	id := c.Param("id")
	idm := c.Param("idm")
	data, err := paketmodels.GetByMembership(id, idm)
	if err != nil {
		c.JSON(406, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"data":   data,
			"status": "ok",
		})
	}
}

func (P *PaketControll) Update(c *gin.Context) {
	id := c.Param("id")
	var data forms.Paket
	if c.BindJSON(&data) != nil {
		c.JSON(405, gin.H{
			"error": "error binding json",
		})
	} else {
		err := paketmodels.Update(id, data)
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

func (P *PaketControll) Updateproduct(c *gin.Context) {
	id := c.Param("id")
	var data struct {
		product []string `json:"product"`
	}
	if c.BindJSON(&data) != nil {
		c.JSON(405, gin.H{"error": "error binding json"})
	} else {
		err := paketmodels.UpdateProduct(id, data.product)
		if err != nil {
			c.JSON(406, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "succes updated",
				"status":  "ok",
			})
		}
	}
}

func (P *PaketControll) Delete(c *gin.Context) {
	id := c.Param("id")
	err := paketmodels.Delete(id)
	if err != nil {
		c.JSON(405, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"message": "deleted",
			"status":  "ok",
		})
	}
}

func (P *PaketControll) ListByMembership(c *gin.Context) {
	membership := c.Query("membership")
	sort := c.Query("sort")
	pageNo := c.Query("page")
	perPage := c.Query("per_page")
	filter := c.Query("filter")
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

	data, err := paketmodels.ListByMembership(membership, filter, sort, pageNo, perPage)
	lastPage := float64(len(data)) / float64(pp)
	if pp != 0 {
		if len(data)%pp == 0 {
			lastPage = lastPage
		} else {
			lastPage += 1
		}
	} else {
		lastPage = float64(len(data)) / float64(5)
	}
	if err != nil {
		c.JSON(405, gin.H{
			"message": "terjadi kesalahan",
			"error":   err.Error(),
		})
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
	}
}
