package Router

import (
	"fmt"
	"simulator/Controllers"
	"simulator/Global"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	gin.SetMode(Global.Mode)
	router := gin.Default()
	api := router.Group("device/api")
	{
		api.Group("check")
		{
			api.HEAD("check", Controllers.Api_Check)
		}

		api.Group("asrs/mission")
		{
			api.POST("asrs/mission", Controllers.Mission)
			api.GET("asrs/mission", Controllers.MissionQuantity)
		}

		api.Group("asrs/status")
		{
			api.POST("asrs/status", Controllers.Status)
		}

		api.Group("asrs/private/mission")
		{
			api.PUT("asrs/private/mission", Controllers.MissionPriveteControl)
		}

		api.Group("asrs/command")
		{
			api.POST("asrs/command", Controllers.SendCommand)
		}
	}
	router.Run(fmt.Sprintf(":%s", Global.Port))
}
