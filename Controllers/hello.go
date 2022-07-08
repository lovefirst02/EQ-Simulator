package Controllers

import "github.com/gin-gonic/gin"

func TestHello(c *gin.Context) {
	c.String(200, "Hello")
}
