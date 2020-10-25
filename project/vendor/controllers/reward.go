package controllers

import (
	"forms"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RewardControll struct{}

func (R *RewardControll) Create(c *gin.Context) {
	var data forms.Rewards
	if c.BindJSON(&data) != nil {
		c.JSON(406, gin.H{
			"error": "binding json",
		})
	} else {
		err := rewardmodels.Create(data)
		if err != nil {
			c.JSON(405, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "Reward Created",
				"status":  "ok",
			})
		}
	}
}

func (R *RewardControll) Get(c *gin.Context) {
	id := c.Param("id")
	data, err := rewardmodels.Get(id)
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

func (R *RewardControll) Delete(c *gin.Context) {
	id := c.Param("id")
	err := rewardmodels.Delete(id)
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

func (R *RewardControll) Update(c *gin.Context) {
	var data forms.Rewards
	if c.BindJSON(&data) != nil {
		c.JSON(405, gin.H{
			"error": "error binding json",
		})
	} else {
		id := c.Param("id")
		err := rewardmodels.Update(id, data)
		if err != nil {
			c.JSON(406, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "updated",
			})
		}

	}
}

func (R *RewardControll) List(c *gin.Context) {
	sort := c.Query("sort")
	pageNo, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	filter := c.Query("filter")
	if sort == "" {
		sort = "id"
	}
	if pageNo == 0 {
		pageNo = 1
	}
	if perPage == 0 {
		perPage = 5
	}
	data, count, err := rewardmodels.List(filter, sort, pageNo, perPage)
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

func (R *RewardControll) ClaimReward(c *gin.Context) {
	account := c.Param("account")
	reward := c.Param("reward")
	err := accountmodels.ClaimReward(account, reward)
	if err != nil {
		c.JSON(405, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"message": "claimmed",
			"status":  "ok",
		})
	}
}
