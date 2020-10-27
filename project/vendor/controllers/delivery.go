package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeliveryControll struct{}

func (D *DeliveryControll) ListProvince(c *gin.Context) {
	sort := c.Query("sort")
	pageNo, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	filter := c.Query("filter")
	if sort == "" {
		sort = "province_id"
	}
	if pageNo == 0 {
		pageNo = 1
	}
	if perPage == 0 {
		perPage = 5
	}
	// pp, _ := perPage)
	// pn, _ := strconv.Atoi(pageNo)

	data, count, err := deliverymodels.GetListProvince(filter, sort, pageNo, perPage)
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

func (D *DeliveryControll) ListCityByProvince(c *gin.Context) {
	sort := c.Query("sort")
	pageNo, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	filter := c.Query("filter")
	prov, _ := strconv.Atoi(c.Param("id"))
	if sort == "" {
		sort = "city_id"
	}
	if pageNo == 0 {
		pageNo = 1
	}
	if perPage == 0 {
		perPage = 5
	}
	// pp, _ := perPage)
	// pn, _ := strconv.Atoi(pageNo)

	data, count, err := deliverymodels.GetListCityByPorvince(prov, filter, sort, pageNo, perPage)
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

func (D *DeliveryControll) ListCity(c *gin.Context) {
	sort := c.Query("sort")
	pageNo, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	filter := c.Query("filter")
	if sort == "" {
		sort = "city_id"
	}
	if pageNo == 0 {
		pageNo = 1
	}
	if perPage == 0 {
		perPage = 5
	}
	// pp, _ := perPage)
	// pn, _ := strconv.Atoi(pageNo)

	data, count, err := deliverymodels.GetListCity(filter, sort, pageNo, perPage)
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

func (D *DeliveryControll) CheckOngkir(c *gin.Context) {
	origin := c.Query("origin")
	destination := c.Query("destination")
	weight := c.Query("weight")
	kurir := c.Query("account")
	data := deliverymodels.CheckOngkirCourir(origin, destination, weight, kurir)
	c.JSON(200, gin.H{
		"data":   data,
		"status": "ok",
	})
}
