package Controllers

import (
	"net/http"
	"simulator/Global"
	"simulator/Models"

	"github.com/gin-gonic/gin"
)

func SendCommand(c *gin.Context) {
	var command Models.Control

	err := c.ShouldBindJSON(&command)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Code":        1,
			"Description": err.Error(),
		})
		return
	}
	if Asrs, ok := Global.Asrs[command.AsrsID]; ok {
		Asrs.AsrsControl(command)
	}

	c.JSON(http.StatusOK, gin.H{
		"Code":        0,
		"Description": "傳送指令成功",
	})
}
