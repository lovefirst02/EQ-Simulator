package Controllers

import (
	"net/http"
	"simulator/Global"
	"simulator/Models"

	"github.com/gin-gonic/gin"
)

func MissionQuantity(c *gin.Context) {
	Mission := make(map[string][]Models.Mission)
	for k, v := range Global.Asrs {
		Mission[k] = v.Mission
	}
	c.JSON(http.StatusOK, gin.H{
		"Code":        0,
		"Description": Mission,
	})
}

func Mission(c *gin.Context) {
	var Mission Models.Mission

	err := c.ShouldBindJSON(&Mission)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Code":        1,
			"Description": err.Error(),
		})
		return
	}
	if Asrs, ok := Global.Asrs[Mission.AsrsID]; ok {
		Asrs.AsrsMission(Mission)
	}

	c.JSON(http.StatusOK, gin.H{
		"Code":        0,
		"Description": "創建任務成功",
	})
}
