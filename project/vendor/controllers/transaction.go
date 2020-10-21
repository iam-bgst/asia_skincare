package controllers

import (
	"fmt"
	"forms"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransactionControll struct{}

func (T *TransactionControll) Add(c *gin.Context) {
	var data forms.Transaction
	if c.BindJSON(&data) != nil {
		c.JSON(405, gin.H{
			"error": "error binding json",
		})
	} else {
		da, err := transactionmodels.Create(data)
		if err != nil {
			c.JSON(405, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "Transaction Created",
				"data":    da,
			})
		}
	}
}

func (T *TransactionControll) AddPicturePay(c *gin.Context) {
	var data forms.Evidence
	if c.BindJSON(&data) != nil {
		c.JSON(405, gin.H{
			"error": "error binding json",
		})
	}
	id := c.Param("id")
	err := transactionmodels.AddPicturePay(id, data)
	if err != nil {
		c.JSON(405, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "add picture success",
		})
	}
}

func (T *TransactionControll) ListTransactionOnagent(c *gin.Context) {
	id_account := c.Param("account")
	sort := c.Query("sort")
	pageNo, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	filter := c.Query("filter")

	status, _ := strconv.Atoi(c.Query("status"))

	if sort == "" {
		sort = "id"
	}
	if pageNo == 0 {
		pageNo = 1
	}
	if perPage == 0 {
		perPage = 5
	}
	data, count, err := transactionmodels.TransactionOnAgent(id_account, filter, sort, pageNo, perPage, status)
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
	if id_account == "" {
		c.JSON(400, gin.H{
			"error": "account not set",
		})
		c.Abort()
	}
	if err != nil {
		c.JSON(400, gin.H{
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

func (T *TransactionControll) ListHistory(c *gin.Context) {
	id_account := c.Param("account")
	sort := c.Query("sort")
	pageNo, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	filter := c.Query("filter")

	status, _ := strconv.Atoi(c.Query("status"))

	if sort == "" {
		sort = "id"
	}
	if pageNo == 0 {
		pageNo = 1
	}
	if perPage == 0 {
		perPage = 5
	}
	data, count, err := transactionmodels.HistoyTransaction(id_account, filter, sort, pageNo, perPage, status)
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
	if id_account == "" {
		c.JSON(400, gin.H{
			"error": "account not set",
		})
		c.Abort()
	}
	if err != nil {
		c.JSON(400, gin.H{
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

func (T *TransactionControll) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	status_code, _ := strconv.Atoi(c.Query("status"))
	err := transactionmodels.UpdateStatus(id, status_code)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"message": fmt.Sprintf("update to %s", transactionmodels.GetStatus(status_code)),
			"status":  "ok",
		})
	}
}

func (T *TransactionControll) AddResiToTransaction(c *gin.Context) {
	id := c.Param("id")
	var AddResi struct {
		Resi string `json:"resi"`
	}
	if c.BindJSON(&AddResi) != nil {
		c.JSON(406, gin.H{
			"error": "error binding json",
		})
	} else {
		err := transactionmodels.UpdateResi(id, AddResi.Resi)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "resi added",
				"status":  "ok",
			})
		}
	}
}
