package middleware

import (
	"strings"
	"swagger"

	"github.com/gin-gonic/gin"
)

func init() {
	r := router.Routes()
	for _, rr := range r {
		swagger.Datadoc.Paths[rr.Path] = gin.H{
			strings.ToLower(rr.Method): gin.H{
				"consumes": []string{
					"application/json",
				},
				"parameters": []gin.H{},
				"responses":  gin.H{},
			},
		}
	}
}
