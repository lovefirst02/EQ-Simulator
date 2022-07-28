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
		api.HEAD("check", Controllers.Api_Check)
	}

	asrs := api.Group("asrs")
	{
		asrs.POST("mission", Controllers.Mission)
		asrs.GET("mission", Controllers.MissionQuantity)
		asrs.POST("status", Controllers.Status)
		asrs.PUT("private/mission", Controllers.MissionPriveteControl)
		asrs.POST("command", Controllers.SendCommand)
	}

	lifter := api.Group("lifter")
	{
		lifter.POST("mission", Controllers.LifterMission)
		lifter.GET("mission", Controllers.LifterMissionQuantity)
		lifter.POST("status", Controllers.LifterStatus)
		lifter.PUT("private/mission", Controllers.LifterMissionPriveteControl)
		lifter.POST("command", Controllers.LifterSendCommand)

	}

	erack := api.Group("erack")
	{
		erack.GET("storage", Controllers.ErackStorageDetail)
		erack.POST("install", Controllers.ErackInstall)
		erack.POST("uninstall", Controllers.ErackUnInstall)
		erack.POST("pre", Controllers.ErackPre)
	}
	router.Run(fmt.Sprintf(":%s", Global.Port))
}
