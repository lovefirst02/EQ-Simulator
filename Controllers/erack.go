package Controllers

import (
	"net/http"
	"reflect"
	"simulator/Global"
	"simulator/Models"
	"simulator/Simulator"
	"simulator/Util"

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
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"Code":        1,
			"Description": "Not Found Erack",
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
		err := Erack.Uninstall(Detail.Storage)
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
		err := Erack.PreStorage(Detail.Storage)
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
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"Code":        1,
			"Description": "Not Found Erack",
		})
	}
}

func GetErackEmptyStorage(c *gin.Context) {
	var Detail map[string]interface{}

	err := c.Bind(&Detail)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Code":        1,
			"Description": err.Error(),
		})
		return
	}

	key := reflect.ValueOf(Detail).MapKeys()
	if !Util.Contain(key, "ErackID") {
		c.JSON(http.StatusBadRequest, gin.H{
			"Code":        1,
			"Description": "No Object Key ErackID",
		})
		return
	}

	if Erack, ok := Global.Erack[Detail["ErackID"].(string)]; ok {
		result := Erack.CheckEmptyStorage()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Code":        1,
				"Description": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"Code":        0,
			"Description": result,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"Code":        1,
			"Description": "Not Found Erack",
		})
	}
}
