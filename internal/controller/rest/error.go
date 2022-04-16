package rest

import (
	"github.com/gin-gonic/gin"
)

func response(c *gin.Context, status int, data interface{}) {
	if data == nil {
		c.Status(status)
	}

	c.JSON(status, gin.H{
		"ok":   true,
		"data": data,
	})

}

func errorResponse(c *gin.Context, status int, msg string) {
	c.AbortWithStatusJSON(status, gin.H{
		"ok":        false,
		"error_msg": msg,
	})
}
