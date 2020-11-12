package controllers

import (
	"forms"

	"github.com/gin-gonic/gin"
)

type MembershipControll struct{}

func (MS *MembershipControll) Create(c *gin.Context) {
	var forms forms.Membership
	if c.BindJSON(&forms) != nil {
		c.JSON(406, gin.H{"error": "error when binding json"})
	} else {
		err := membershipmodels.Add(forms.Name)
		if err != nil {
			c.JSON(405, gin.H{"error": "error while add membership"})
		} else {
			c.JSON(200, gin.H{
				"message": "Membership Added",
				"status":  "ok",
			})
		}
	}
}

func (MS *MembershipControll) ListAll(c *gin.Context) {
	ne := c.Query("ne")
	code := c.Query("code")
	data, err := membershipmodels.ListAll(ne, code)

	if err != nil {
		c.JSON(405, gin.H{
			"error": "error while get list membership",
		})
	} else {
		c.JSON(200, gin.H{
			"data":   data,
			"status": "ok",
		})
	}
}
