package Router

import (
	"simulator/Controllers"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	api := router.Group("device/api")
	{
		api.HEAD("check", Controllers.Api_Check)
		api.POST("asrs/mission", Controllers.Mission)
		api.GET("asrs/mission", Controllers.MissionQuantity)
		api.POST("asrs/status", Controllers.Status)
		api.POST("asrs/command", Controllers.SendCommand)
		api.GET("asrs/hello", Controllers.TestHello)
	}
	router.Run(":8880")
}
