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
	}

	asrs := api.Group("asrs")
	{
		asrs.Group("mission")
		{
			asrs.POST("mission", Controllers.Mission)
			asrs.GET("mission", Controllers.MissionQuantity)
		}
		asrs.Group("status")
		{
			asrs.POST("status", Controllers.Status)
		}
		asrs.Group("/private/mission")
		{
			asrs.PUT("private/mission", Controllers.MissionPriveteControl)
		}
		asrs.Group("command")
		{
			asrs.POST("command", Controllers.SendCommand)
		}
	}

	lifter := api.Group("lifter")
	{
		lifter.Group("mission")
		{
			lifter.POST("mission", Controllers.LifterMission)
			lifter.GET("mission", Controllers.LifterMissionQuantity)
		}
		lifter.Group("status")
		{
			lifter.POST("status", Controllers.LifterStatus)
		}
		lifter.Group("/private/mission")
		{
			lifter.PUT("private/mission", Controllers.LifterMissionPriveteControl)
		}
		lifter.Group("command")
		{
			lifter.POST("command", Controllers.LifterSendCommand)
		}
	}
	router.Run(fmt.Sprintf(":%s", Global.Port))
}
