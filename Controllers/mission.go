package Controllers

import (
	"net/http"
	"simulator/Global"
	"simulator/Models"

	"github.com/gin-gonic/gin"
)

func MissionPriveteControl(c *gin.Context) {

	var Control Models.MissionPrivateControl

	err := c.ShouldBindJSON(&Control)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Code":        1,
			"Description": err.Error(),
		})
		return
	}

	if v, ok := Global.Asrs[Control.AsrsID]; ok {
		result := v.AsrsMissionPrivateControl(Control)
		if result {
			c.JSON(http.StatusOK, gin.H{
				"Code":        0,
				"Description": "任務指令傳送成功",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"Code":        1,
				"Description": "任務指令傳送失敗",
			})

		}
	}

}

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
		Mission.Control = make(chan string)
		Asrs.AsrsMission(Mission)
	}

	c.JSON(http.StatusOK, gin.H{
		"Code":        0,
		"Description": "創建任務成功",
	})
}

/////////////////////////////LIFTER///////////////////////////////
func LifterMissionPriveteControl(c *gin.Context) {

	var Control Models.LifterMissionPrivateControl

	err := c.ShouldBindJSON(&Control)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Code":        1,
			"Description": err.Error(),
		})
		return
	}

	if v, ok := Global.Lifter[Control.LifterID]; ok {
		result := v.LifterMissionPrivateControl(Control)
		if result {
			c.JSON(http.StatusOK, gin.H{
				"Code":        0,
				"Description": "任務指令傳送成功",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"Code":        1,
				"Description": "任務指令傳送失敗",
			})

		}
	}

}

func LifterMissionQuantity(c *gin.Context) {
	Mission := make(map[string][]Models.LifterMission)
	for k, v := range Global.Lifter {
		Mission[k] = v.Mission
	}

	c.JSON(http.StatusOK, gin.H{
		"Code":        0,
		"Description": Mission,
	})
}

func LifterMission(c *gin.Context) {
	var Mission Models.LifterMission

	err := c.ShouldBindJSON(&Mission)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Code":        1,
			"Description": err.Error(),
		})
		return
	}
	if Lifter, ok := Global.Lifter[Mission.LifterID]; ok {
		Mission.Control = make(chan string)
		Lifter.LifterMission(Mission)
	}

	c.JSON(http.StatusOK, gin.H{
		"Code":        0,
		"Description": "創建任務成功",
	})
}
