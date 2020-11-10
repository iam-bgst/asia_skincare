package controllers

import (
	"forms"
	"models"

	"github.com/gin-gonic/gin"
)

type HeaderControll struct{}

func (H *HeaderControll) Create(c *gin.Context) {
	var data forms.Header
	if c.BindJSON(&data) != nil {
		c.JSON(405, gin.H{
			"error": "error binding json",
		})
	}
}

func (H *HeaderControll) Get(c *gin.Context) {
	data, err := headermodels.Get()
	if err != nil {
		c.JSON(406, gin.H{
			"err": err.Error(),
		})
	} else {
		if data == (models.Header{}) {
			c.JSON(200, gin.H{
				"data": gin.H{
					"id":    "",
					"title": "Belum ada title",
				},
				"status": "ok",
			})
		} else {
			c.JSON(200, gin.H{
				"data": data,
			})
		}
	}
}

func (H *HeaderControll) Update(c *gin.Context) {
	var data forms.Header
	if c.BindJSON(&data) != nil {
		c.JSON(405, gin.H{
			"error": "error binding json",
		})
	} else {
		id := c.Param("id")
		err := headermodels.Update(id, data)
		if err != nil {
			c.JSON(406, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "Updated",
			})
		}
	}
}
