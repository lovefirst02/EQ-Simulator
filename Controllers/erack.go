package Controllers

import (
	"net/http"
	"simulator/Global"
	"simulator/Models"
	"simulator/Simulator"

	"github.com/gin-gonic/gin"
)

func ErackStorageDetail(c *gin.Context) {
	Detail := make(map[string]map[string]Simulator.StorageDetail)
	for k, v := range Global.Erack {
		Detail[k] = v.Storage
	}
	c.JSON(http.StatusOK, gin.H{
		"Code":        0,
		"Description": Detail,
	})
}

func ErackInstall(c *gin.Context) {
	var Detail Models.ErackStorage

	err := c.ShouldBindJSON(&Detail)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Code":        1,
			"Description": err.Error(),
		})
		return
	}
	if Erack, ok := Global.Erack[Detail.ErackID]; ok {
		err := Erack.Install(Detail.Storage, Detail.CarrierID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Code":        1,
				"Description": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"Code":        0,
			"Description": "成功",
		})
	}
}

func ErackUnInstall(c *gin.Context) {
	var Detail Models.ErackStorage

	err := c.ShouldBindJSON(&Detail)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Code":        1,
			"Description": err.Error(),
		})
	}
	if Erack, ok := Global.Erack[Detail.ErackID]; ok {
		Erack.Uninstall(Detail.Storage)
		c.JSON(http.StatusOK, gin.H{
			"Code":        0,
			"Description": "成功",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"Code":        1,
			"Description": "Not Found Erack",
		})
	}
}

func ErackPre(c *gin.Context) {
	var Detail Models.ErackStorage

	err := c.ShouldBindJSON(&Detail)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Code":        1,
			"Description": err.Error(),
		})
	}
	if Erack, ok := Global.Erack[Detail.ErackID]; ok {
		Erack.PreStorage(Detail.Storage)
		c.JSON(http.StatusOK, gin.H{
			"Code":        0,
			"Description": "成功",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"Code":        1,
			"Description": "Not Found Erack",
		})
	}
}
