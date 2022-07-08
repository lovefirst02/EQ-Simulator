package Controllers

import "github.com/gin-gonic/gin"

func Api_Check(c *gin.Context) {
	c.JSON(200, gin.H{
		"Code":        0,
		"Description": "Connect",
	})
}
